package charting

import (
	"fmt"
	"image/color"
	"image/png"
	"io"
	"os"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"tmobile-stats/internal/analysis"
	"tmobile-stats/internal/models"
)

// logTicks is a custom ticker for logarithmic scales that formats labels as integers.
type logTicks struct{}

func (logTicks) Ticks(min, max float64) []plot.Tick {
	ticks := plot.LogTicks{}.Ticks(min, max)
	for i := range ticks {
		if ticks[i].Label != "" {
			// Format as standard decimal/integer
			ticks[i].Label = fmt.Sprintf("%g", ticks[i].Value)
		}
	}
	return ticks
}

// downsample reduces the resolution of the data by averaging points into buckets.
func downsample(data plotter.XYs, maxPoints int) plotter.XYs {
	if len(data) <= maxPoints {
		return data
	}

	bucketSize := float64(len(data)) / float64(maxPoints)
	downsampled := make(plotter.XYs, 0, maxPoints)

	for i := 0; i < maxPoints; i++ {
		start := int(float64(i) * bucketSize)
		end := int(float64(i+1) * bucketSize)
		if end > len(data) {
			end = len(data)
		}
		if start >= end {
			continue
		}

		var sumY float64
		var sumX float64
		count := 0.0

		for k := start; k < end; k++ {
			sumX += data[k].X
			sumY += data[k].Y
			count++
		}

		if count > 0 {
			downsampled = append(downsampled, plotter.XY{
				X: sumX / count,
				Y: sumY / count,
			})
		}
	}

	return downsampled
}

// Generate creates a PNG chart from the provided stats and saves it to outputFile.
func Generate(data []models.CombinedStats, outputFile string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return GenerateToWriter(data, f)
}

// GenerateToWriter creates a PNG chart from the provided stats and writes it to the provided writer.
func GenerateToWriter(data []models.CombinedStats, w io.Writer) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to chart")
	}

	// Prepare plotters
	latencyXYs := make(plotter.XYs, len(data))
	stdDevXYs := make(plotter.XYs, len(data))
	rsrpXYs := make(plotter.XYs, len(data))
	sinrXYs := make(plotter.XYs, len(data))
	lossXYs := make(plotter.XYs, len(data))

	// Helper to safely cast time
	getTime := func(t int64) float64 {
		return float64(t)
	}

	for i, d := range data {
		t := getTime(d.Gateway.Time.LocalTime)

		// Sanitize for Log Scale (must be > 0)
		avg := d.Ping.Avg
		if avg <= 0 {
			avg = 0.1
		}
		sd := d.Ping.StdDev
		if sd <= 0 {
			sd = 0.1
		}
		loss := d.Ping.Loss
		if loss <= 0 {
			loss = 0.1
		}

		latencyXYs[i].X = t
		latencyXYs[i].Y = avg

		stdDevXYs[i].X = t
		stdDevXYs[i].Y = sd

		rsrpXYs[i].X = t
		rsrpXYs[i].Y = float64(d.Gateway.Signal.FiveG.RSRP)

		sinrXYs[i].X = t
		sinrXYs[i].Y = float64(d.Gateway.Signal.FiveG.SINR)

		lossXYs[i].X = t
		lossXYs[i].Y = loss
	}

	// Common Time Ticker
	// We use the standard TimeTicks which interprets float64 X values as unix seconds
	timeTicks := plot.TimeTicks{
		Format: "15:04",
		Time: func(t float64) time.Time {
			return time.Unix(int64(t), 0).Local()
		},
	}

	// Create Plots

	// 1. Latency & Packet Loss
	pLat := plot.New()
	pLat.Title.Text = "Latency & Packet Loss"
	pLat.Y.Label.Text = "ms / %"
	pLat.X.Tick.Marker = timeTicks
	pLat.Y.Scale = plot.LogScale{}
	pLat.Y.Tick.Marker = logTicks{}

	lineLat, _ := plotter.NewLine(latencyXYs)
	lineLat.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue

	lineStd, _ := plotter.NewLine(stdDevXYs)
	lineStd.Color = color.RGBA{R: 255, G: 140, B: 0, A: 255} // Dark Orange
	// Solid line (default)

	lineLoss, _ := plotter.NewLine(lossXYs)
	lineLoss.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black

	pLat.Add(lineLat, lineStd, lineLoss)
	pLat.Legend.Add("Avg (ms)", lineLat)
	pLat.Legend.Add("StdDev", lineStd)
	pLat.Legend.Add("Loss (%)", lineLoss)
	pLat.Add(plotter.NewGrid())

	// 2. Signal Strength (Split)
	// 2a. SINR (Top)
	pSINR := plot.New()
	pSINR.Title.Text = "Signal Strength (SINR / RSRP)"
	pSINR.Y.Label.Text = "SINR (dB)"
	pSINR.X.Tick.Marker = timeTicks
	pSINR.X.Tick.Label.Color = color.Transparent // Hide labels, keep spacing

	lineSINR, _ := plotter.NewLine(sinrXYs)
	lineSINR.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green

	pSINR.Add(lineSINR)
	pSINR.Legend.Add("SINR", lineSINR)
	pSINR.Add(plotter.NewGrid())

	// 2b. RSRP (Bottom)
	pRSRP := plot.New()
	// No title for bottom plot
	pRSRP.Y.Label.Text = "RSRP (dBm)"
	pRSRP.X.Tick.Marker = timeTicks

	lineRSRP, _ := plotter.NewLine(rsrpXYs)
	lineRSRP.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red

	pRSRP.Add(lineRSRP)
	pRSRP.Legend.Add("RSRP", lineRSRP)
	pRSRP.Add(plotter.NewGrid())

	bandXYs := make(plotter.XYs, len(data))
	towerXYs := make(plotter.XYs, len(data))

	// Collect Unique Towers for Mapping
	towerSet := make(map[int]bool)
	for _, d := range data {
		if d.Gateway.Signal.FiveG.GNBID > 0 {
			towerSet[d.Gateway.Signal.FiveG.GNBID] = true
		}
	}
	var sortedTowers []int
	for t := range towerSet {
		sortedTowers = append(sortedTowers, t)
	}
	
	// Simple Bubble Sort
	for i := 0; i < len(sortedTowers); i++ {
		for j := i + 1; j < len(sortedTowers); j++ {
			if sortedTowers[i] > sortedTowers[j] {
				sortedTowers[i], sortedTowers[j] = sortedTowers[j], sortedTowers[i]
			}
		}
	}

	getTowerY := func(gnbid int) float64 {
		if gnbid == 0 {
			return 0
		}
		for i, t := range sortedTowers {
			if t == gnbid {
				// Bands are 1, 2, 3. Leave 4 empty. Start Towers at 5.
				return float64(i + 5)
			}
		}
		return 0
	}

	for i, d := range data {
		t := getTime(d.Gateway.Time.LocalTime)
		bandXYs[i].X = t
		towerXYs[i].X = t

		// Map bands to levels
		// Priority: n41 > n25 > n71
		level := 0.0 // No signal/Unknown
		hasBand := func(target string) bool {
			for _, b := range d.Gateway.Signal.FiveG.Bands {
				if b == target {
					return true
				}
			}
			return false
		}

		if hasBand("n41") {
			level = 3
		} else if hasBand("n25") {
			level = 2
		} else if hasBand("n71") {
			level = 1
		}
		bandXYs[i].Y = level

		// Map Towers
		towerXYs[i].Y = getTowerY(d.Gateway.Signal.FiveG.GNBID)
	}

	// Apply Smoothing to Towers if "All" smoothing is active
	// (Same trigger as Bands)
	if len(data) > 1 {
		startTime := data[0].Gateway.Time.LocalTime
		endTime := data[len(data)-1].Gateway.Time.LocalTime
		duration := float64(endTime - startTime)
		
		if duration > 86400 && len(towerXYs) > 600 {
			towerXYs = downsample(towerXYs, 600)
		}
	}

	lineBand, _ := plotter.NewLine(bandXYs)
	lineBand.StepStyle = plotter.PreStep
	lineBand.Color = color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange

	lineTower, _ := plotter.NewLine(towerXYs)
	lineTower.StepStyle = plotter.PreStep
	lineTower.Color = color.RGBA{R: 255, G: 0, B: 255, A: 255} // Magenta

	// 3. 5G Band & Tower Plot
	pBand := plot.New()
	pBand.Title.Text = "Connection Info / Tower & Radio Bands"
	pBand.X.Tick.Marker = timeTicks

	// Custom Y Ticks for Bands & Towers
	// Min = 0, Max = 5 + len(towers)
	pBand.Y.Min = 0
	pBand.Y.Max = float64(5 + len(sortedTowers))
	
	ticks := []plot.Tick{
		{Value: 1, Label: "n71"},
		{Value: 2, Label: "n25"},
		{Value: 3, Label: "n41"},
	}
	for i, tID := range sortedTowers {
		ticks = append(ticks, plot.Tick{
			Value: float64(i + 5),
			Label: fmt.Sprintf("%d", tID),
		})
	}
	pBand.Y.Tick.Marker = plot.ConstantTicks(ticks)

	pBand.Add(lineBand, lineTower)
	pBand.Add(plotter.NewGrid())

	// 4. Signal Bars Plot (New)
	pBars := plot.New()
	pBars.Title.Text = "Signal Bars"
	pBars.Y.Label.Text = "Bars"
	pBars.X.Tick.Marker = timeTicks
	pBars.Y.Min = 0
	pBars.Y.Max = 5
	pBars.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: 1, Label: "1"},
		{Value: 2, Label: "2"},
		{Value: 3, Label: "3"},
		{Value: 4, Label: "4"},
	})

	barsXYs := make(plotter.XYs, len(data))
	healthXYs := make(plotter.XYs, len(data))
	for i, d := range data {
		t := getTime(d.Gateway.Time.LocalTime)
		barsXYs[i].X = t
		barsXYs[i].Y = float64(d.Gateway.Signal.FiveG.Bars)

		healthXYs[i].X = t
		healthXYs[i].Y = analysis.CalculateSignalHealth(d.Gateway.Signal.FiveG.RSRP, d.Gateway.Signal.FiveG.SINR)
	}

	// Check if smoothing is needed
	// Criteria: Duration > 2 hours AND Data Points > 300
	shouldSmoothBars := false

	if len(data) > 1 {
		startTime := data[0].Gateway.Time.LocalTime
		endTime := data[len(data)-1].Gateway.Time.LocalTime
		duration := float64(endTime - startTime)
		
		// 2 hours = 7200 seconds
		if duration > 7200 && len(barsXYs) > 300 {
			shouldSmoothBars = true
			// Downsample Bars to target ~300 points
			barsXYs = downsample(barsXYs, 300)
			healthXYs = downsample(healthXYs, 300)
		}

		// 24 hours = 86400 seconds
		if duration > 86400 && len(latencyXYs) > 600 {
			// Downsample others to target ~600 points (higher fidelity)
			latencyXYs = downsample(latencyXYs, 600)
			stdDevXYs = downsample(stdDevXYs, 600)
			lossXYs = downsample(lossXYs, 600)
			rsrpXYs = downsample(rsrpXYs, 600)
			sinrXYs = downsample(sinrXYs, 600)
			bandXYs = downsample(bandXYs, 600)
		}
	}

	// 1. Health Area (Background)
	// Create a polygon for the filled area (requires closing the path to Y=0)
	healthAreaXYs := make(plotter.XYs, len(healthXYs)+2)
	copy(healthAreaXYs[1:], healthXYs)
	healthAreaXYs[0] = plotter.XY{X: healthXYs[0].X, Y: 0}
	healthAreaXYs[len(healthAreaXYs)-1] = plotter.XY{X: healthXYs[len(healthXYs)-1].X, Y: 0}

	polyHealth, _ := plotter.NewPolygon(healthAreaXYs)
	polyHealth.Color = color.RGBA{R: 240, G: 240, B: 240, A: 255} // Very Light Grey
	polyHealth.LineStyle.Width = 0                                // No border for the area

	// 2. Health Line (Top of area)
	lineHealth, _ := plotter.NewLine(healthXYs)
	lineHealth.Color = color.RGBA{R: 169, G: 169, B: 169, A: 255} // Dark Grey border
	lineHealth.Width = vg.Points(1)

	// 3. Reported Bars (Foreground)
	lineBars, _ := plotter.NewLine(barsXYs)
	lineBars.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black
	
	if shouldSmoothBars {
		// Smoothed Mode: Normal Line (Continuous), Thinner
		lineBars.StepStyle = plotter.NoStep
		lineBars.Width = vg.Points(0.8) // Keep it thin
	} else {
		// Raw Mode: Step Line (Discrete)
		lineBars.StepStyle = plotter.PreStep
		lineBars.Width = vg.Points(0.8)
	}

	pBars.Add(polyHealth, lineHealth, lineBars)
	pBars.Legend.Add("Reported Bars", lineBars)
	pBars.Legend.Add("Signal Health", lineHealth)
	pBars.Add(plotter.NewGrid())

	// Combine into a single image
	const width = 20 * vg.Inch  // Double width
	const height = 10 * vg.Inch // Adjusted height for 2 rows

	c := vgimg.NewWith(vgimg.UseWH(width, height), vgimg.UseBackgroundColor(color.White))
	dc := draw.New(c)

	// Layout: 2x2 Grid
	// Col 1: Left, Col 2: Right
	// Row 1: Top, Row 2: Bottom

	colWidth := width / 2
	rowHeight := height / 2

	// 1. Latency (Top Left)
	rectLat := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: rowHeight},
			Max: vg.Point{X: colWidth, Y: height},
		},
	}
	pLat.Draw(rectLat)

	// 2. Signal Strength (Top Right - Split)
	// SINR (Top half of Top Right)
	// Quadrant Height = rowHeight. Split at rowHeight + (rowHeight/2).
	rectSINR := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: colWidth, Y: rowHeight + (rowHeight / 2)},
			Max: vg.Point{X: width, Y: height},
		},
	}
	pSINR.Draw(rectSINR)

	// RSRP (Bottom half of Top Right)
	rectRSRP := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: colWidth, Y: rowHeight},
			Max: vg.Point{X: width, Y: rowHeight + (rowHeight / 2)},
		},
	}
	pRSRP.Draw(rectRSRP)

	// 3. Bars (Bottom Left)
	rectBars := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: colWidth, Y: rowHeight},
		},
	}
	pBars.Draw(rectBars)

	// 4. Bands (Bottom Right)
	rectBand := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: colWidth, Y: 0},
			Max: vg.Point{X: width, Y: rowHeight},
		},
	}
	pBand.Draw(rectBand)

	return png.Encode(w, c.Image())
}