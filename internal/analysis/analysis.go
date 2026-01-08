package analysis

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"tmobile-stats/internal/models"
)

type Metric struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func (m *Metric) Add(val float64) {
	if m.Count == 0 || val < m.Min {
		m.Min = val
	}
	if m.Count == 0 || val > m.Max {
		m.Max = val
	}
	m.Sum += val
	m.Count++
}

func (m Metric) Avg() float64 {
	if m.Count == 0 {
		return 0
	}
	return m.Sum / float64(m.Count)
}

type Report struct {
	TotalSamples int
	StartTime    time.Time
	EndTime      time.Time
	Filter       *TimeFilter

	RSRP Metric
	SINR Metric
	Ping Metric
	Loss Metric

	TotalPingSent int
	TotalPingLost int

	Bands  map[string]int
	Towers map[int]int
	Bars   map[float64]int
	LastTowerID int
	LastBars    float64

	AvgBarsOverall  float64
	AvgBars1h       float64
	AvgSignalHealth float64
	Has1hData       bool
}

func Run(path string, filter *TimeFilter) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return Analyze(file, os.Stdout, filter)
}

func Analyze(input io.Reader, output io.Writer, filter *TimeFilter) error {
	report := &Report{
		Bands:  make(map[string]int),
		Towers: make(map[int]int),
		Bars:   make(map[float64]int),
		Filter: filter,
	}
	// Initialize Min values to avoid 0.0 bias
	report.Ping.Min = math.MaxFloat64
	report.RSRP.Min = 0
	report.SINR.Min = 99

	// Fetch raw data using the new exported parser
	data, err := ParseLog(input, filter)
	if err != nil {
		return err
	}

	var sumBars float64
	var sumHealth float64

	for _, stats := range data {
		report.TotalSamples++
		// Use UpTime as a proxy for relative time if LocalTime isn't enough,
		// but let's assume we can use the current system time if needed,
		// or better, Gateway.Time.LocalTime.
		sampleTime := time.Unix(stats.Gateway.Time.LocalTime, 0)
		if report.StartTime.IsZero() || sampleTime.Before(report.StartTime) {
			report.StartTime = sampleTime
		}
		if sampleTime.After(report.EndTime) {
			report.EndTime = sampleTime
		}

		report.RSRP.Add(float64(stats.Gateway.Signal.FiveG.RSRP))
		report.SINR.Add(float64(stats.Gateway.Signal.FiveG.SINR))

		// Accumulate Bars & Health
		sumBars += stats.Gateway.Signal.FiveG.Bars
		sumHealth += CalculateSignalHealth(stats.Gateway.Signal.FiveG.RSRP, stats.Gateway.Signal.FiveG.SINR)

		if stats.Ping.Received > 0 {
			// Ignore 0.0 pings for Min calculation as it was a bug in earlier versions
			if stats.Ping.Min > 0 {
				if report.Ping.Count == 0 || stats.Ping.Min < report.Ping.Min {
					report.Ping.Min = stats.Ping.Min
				}
			}
			if report.Ping.Count == 0 || stats.Ping.Max > report.Ping.Max {
				report.Ping.Max = stats.Ping.Max
			}
			report.Ping.Sum += stats.Ping.Avg
			report.Ping.Count++
		}
		report.Loss.Add(stats.Ping.Loss)
		report.TotalPingSent += stats.Ping.Sent
		report.TotalPingLost += stats.Ping.Sent - stats.Ping.Received

		for _, b := range stats.Gateway.Signal.FiveG.Bands {
			report.Bands[b]++
		}
		towerID := stats.Gateway.Signal.FiveG.GNBID
		if towerID != 0 {
			report.Towers[towerID]++
			report.LastTowerID = towerID
		}

		report.Bars[stats.Gateway.Signal.FiveG.Bars]++
		report.LastBars = stats.Gateway.Signal.FiveG.Bars
	}

	// Finalize Averages
	if report.TotalSamples > 0 {
		report.AvgBarsOverall = sumBars / float64(report.TotalSamples)
		report.AvgSignalHealth = sumHealth / float64(report.TotalSamples)

		// Calculate Last 1h
		oneHourAgo := report.EndTime.Add(-1 * time.Hour)
		var sumBars1h float64
		var count1h int

		// Iterate backwards from the end of data slice
		for i := len(data) - 1; i >= 0; i-- {
			sampleTime := time.Unix(data[i].Gateway.Time.LocalTime, 0)
			if sampleTime.Before(oneHourAgo) {
				break
			}
			sumBars1h += data[i].Gateway.Signal.FiveG.Bars
			count1h++
		}

		// Only show Last 1h if we have at least 55m of data duration
		if report.EndTime.Sub(report.StartTime) >= 55*time.Minute && count1h > 0 {
			report.AvgBars1h = sumBars1h / float64(count1h)
			report.Has1hData = true
		}
	}

	printReport(output, report)
	return nil
}

// ParseLog reads the provided reader and returns a slice of CombinedStats.
// It skips malformed lines.
func ParseLog(r io.Reader, filter *TimeFilter) ([]models.CombinedStats, error) {
	var results []models.CombinedStats
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var stats models.CombinedStats
		if err := json.Unmarshal(scanner.Bytes(), &stats); err != nil {
			continue // Skip malformed lines
		}

		// Filter by time
		sampleTime := time.Unix(stats.Gateway.Time.LocalTime, 0)
		if filter != nil && !filter.Contains(sampleTime) {
			continue
		}

		results = append(results, stats)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func printReport(w io.Writer, r *Report) {
	fmt.Fprintln(w, "================================================================================")
	fmt.Fprintln(w, " HISTORICAL SIGNAL ANALYSIS")
	fmt.Fprintln(w, "================================================================================")

	if r.TotalSamples == 0 {
		fmt.Fprintln(w, "No data samples found.")
		return
	}

	duration := r.EndTime.Sub(r.StartTime)
	if r.Filter != nil && (!r.Filter.Start.IsZero() || !r.Filter.End.IsZero()) {
		startS := "Begin"
		if !r.Filter.Start.IsZero() {
			startS = r.Filter.Start.Format("2006-01-02 15:04:05")
		}
		endS := "End"
		if !r.Filter.End.IsZero() {
			endS = r.Filter.End.Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(w, "Filter:        %s to %s\n", startS, endS)
	}
	fmt.Fprintf(w, "Data Range:    %s to %s\n", r.StartTime.Format("2006-01-02 15:04:05"), r.EndTime.Format("15:04:05"))
	fmt.Fprintf(w, "Duration:      %v\n", duration.Round(time.Second))
	fmt.Fprintf(w, "Total Samples: %d\n", r.TotalSamples)
	fmt.Fprintln(w)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "METRIC\tMIN\tAVG\tMAX")
	fmt.Fprintln(tw, "------\t---\t---\t---")
	fmt.Fprintf(tw, "RSRP (dBm)\t%.0f\t%.1f\t%.0f\n", r.RSRP.Min, r.RSRP.Avg(), r.RSRP.Max)
	fmt.Fprintf(tw, "SINR (dB)\t%.0f\t%.1f\t%.0f\n", r.SINR.Min, r.SINR.Avg(), r.SINR.Max)
	pMin := r.Ping.Min
	if pMin == math.MaxFloat64 {
		pMin = 0
	}
	fmt.Fprintf(tw, "Ping (ms)\t%.1f\t%.1f\t%.1f\n", pMin, r.Ping.Avg(), r.Ping.Max)
	tw.Flush()

	if r.TotalPingSent > 0 {
		globalLoss := float64(r.TotalPingLost) / float64(r.TotalPingSent) * 100
		fmt.Fprintf(w, "\nRELIABILITY:\n")
		fmt.Fprintf(w, "  Packet Loss: %d / %d (%.2f%%)\n", r.TotalPingLost, r.TotalPingSent, globalLoss)
	}

	fmt.Fprintln(w, "\nBANDS SEEN:")
	printMap(w, r.Bands, r.TotalSamples, duration)

	fmt.Fprintln(w, "\nTOWERS SEEN:")
	printTowerMap(w, r.Towers, r.TotalSamples, r.LastTowerID, duration)

	fmt.Fprintln(w, "\nBARS SEEN:")
	printFloatMap(w, r.Bars, r.TotalSamples, r.LastBars, duration)

	fmt.Fprintln(w, "\nBARS AVG:")
	tw2 := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw2, "Overall\t%.1f\n", r.AvgBarsOverall)
	if r.Has1hData {
		fmt.Fprintf(tw2, "Last 1h\t%.1f\n", r.AvgBars1h)
	}
	fmt.Fprintf(tw2, "SgnlHealth\t%.1f\n", r.AvgSignalHealth)
	tw2.Flush()

	fmt.Fprintln(w, "================================================================================")
}

func printMap(w io.Writer, m map[string]int, total int, totalDuration time.Duration) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, k := range keys {
		count := m[k]
		pct := float64(count) / float64(total) * 100
		dur := time.Duration(float64(totalDuration) * (float64(count) / float64(total)))
		fmt.Fprintf(tw, "  %s\t%d samples (%.1f%%)\t%s\n", k, count, pct, formatSmartDuration(dur))
	}
	tw.Flush()
}

func printTowerMap(w io.Writer, m map[int]int, total int, liveTowerID int, totalDuration time.Duration) {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, k := range keys {
		count := m[k]
		pct := float64(count) / float64(total) * 100
		suffix := ""
		if k == liveTowerID {
			suffix = " live"
		}
		dur := time.Duration(float64(totalDuration) * (float64(count) / float64(total)))
		fmt.Fprintf(tw, "  %d\t%d samples (%.1f%%)%s\t%s\n", k, count, pct, suffix, formatSmartDuration(dur))
	}
	tw.Flush()
}

func printFloatMap(w io.Writer, m map[float64]int, total int, realTimeVal float64, totalDuration time.Duration) {
	keys := make([]float64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, k := range keys {
		count := m[k]
		pct := float64(count) / float64(total) * 100
		suffix := ""
		if k == realTimeVal {
			suffix = " real-time"
		}
		dur := time.Duration(float64(totalDuration) * (float64(count) / float64(total)))
		fmt.Fprintf(tw, "  %g\t%d samples (%.1f%%)%s\t%s\n", k, count, pct, suffix, formatSmartDuration(dur))
	}
	tw.Flush()
}

func formatSmartDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		if m > 0 {
			return fmt.Sprintf("%dh %dm", h, m)
		}
		return fmt.Sprintf("%dh", h)
	}
	if m > 0 {
		if s > 0 && m < 10 { // Only show seconds if minutes are low count
			return fmt.Sprintf("%dm %ds", m, s)
		}
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%ds", s)
}
