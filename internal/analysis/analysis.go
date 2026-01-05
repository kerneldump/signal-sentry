package analysis

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	
	RSRP   Metric
	SINR   Metric
	Ping   Metric
	Loss   Metric
	
	Bands  map[string]int
	Towers map[int]int
}

func Run(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	report := &Report{
		Bands:  make(map[string]int),
		Towers: make(map[int]int),
	}
	// Initialize Min values to avoid 0.0 bias
	report.Ping.Min = 999999
	report.RSRP.Min = 0
	report.SINR.Min = 99

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var stats models.CombinedStats
		if err := json.Unmarshal(scanner.Bytes(), &stats); err != nil {
			continue // Skip malformed lines
		}

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

		for _, b := range stats.Gateway.Signal.FiveG.Bands {
			report.Bands[b]++
		}
		towerID := stats.Gateway.Signal.FiveG.GNBID
		if towerID != 0 {
			report.Towers[towerID]++
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	printReport(report)
	return nil
}

func printReport(r *Report) {
	fmt.Println("================================================================================")
	fmt.Println(" HISTORICAL SIGNAL ANALYSIS")
	fmt.Println("================================================================================")
	
	if r.TotalSamples == 0 {
		fmt.Println("No data samples found.")
		return
	}

	duration := r.EndTime.Sub(r.StartTime)
	fmt.Printf("Time Range:    %s to %s\n", r.StartTime.Format("2006-01-02 15:04:05"), r.EndTime.Format("15:04:05"))
	fmt.Printf("Duration:      %v\n", duration.Round(time.Second))
	fmt.Printf("Total Samples: %d\n", r.TotalSamples)
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC\tMIN\tAVG\tMAX")
	fmt.Fprintln(w, "------\t---\t---\t---")
	fmt.Fprintf(w, "RSRP (dBm)\t%.0f\t%.1f\t%.0f\n", r.RSRP.Min, r.RSRP.Avg(), r.RSRP.Max)
	fmt.Fprintf(w, "SINR (dB)\t%.0f\t%.1f\t%.0f\n", r.SINR.Min, r.SINR.Avg(), r.SINR.Max)
	pMin := r.Ping.Min
	if pMin == 999999 {
		pMin = 0
	}
	fmt.Fprintf(w, "Ping (ms)\t%.1f\t%.1f\t%.1f\n", pMin, r.Ping.Avg(), r.Ping.Max)
	fmt.Fprintf(w, "Loss (%%)\t%.1f\t%.1f\t%.1f\n", r.Loss.Min, r.Loss.Avg(), r.Loss.Max)
	w.Flush()

	fmt.Println("\nBANDS SEEN:")
	printMap(r.Bands, r.TotalSamples)

	fmt.Println("\nTOWERS SEEN:")
	printTowerMap(r.Towers, r.TotalSamples)
	
	fmt.Println("================================================================================")
}

func printMap(m map[string]int, total int) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		pct := float64(m[k]) / float64(total) * 100
		fmt.Printf("  %-10s %d samples (%.1f%%)\n", k, m[k], pct)
	}
}

func printTowerMap(m map[int]int, total int) {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		pct := float64(m[k]) / float64(total) * 100
		fmt.Printf("  %-10d %d samples (%.1f%%)\n", k, m[k], pct)
	}
}
