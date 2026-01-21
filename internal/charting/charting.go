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

// addLastPointLabel adds a text label to the last point of the data series.
func addLastPointLabel(p *plot.Plot, data plotter.XYs, format string, c color.Color) {
	if len(data) == 0 {
		return
	}
	lastPt := data[len(data)-1]

	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    plotter.XYs{lastPt},
		Labels: []string{fmt.Sprintf(format, lastPt.Y)},
	})
	if err != nil {
		return
	}

	labels.TextStyle[0].Color = c
	labels.TextStyle[0].Font.Size = vg.Points(9)
	labels.TextStyle[0].XAlign = draw.XLeft
	labels.TextStyle[0].YAlign = draw.YCenter
	labels.Offset = vg.Point{X: vg.Points(5), Y: 0}

	p.Add(labels)
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

	// ---------------------------
	// 1. Data Preparation
	// ---------------------------

	// Prepare plotters
	latencyXYs := make(plotter.XYs, len(data))
	stdDevXYs := make(plotter.XYs, len(data))
	rsrpXYs := make(plotter.XYs, len(data))
	sinrXYs := make(plotter.XYs, len(data))
	lossXYs := make(plotter.XYs, len(data))
	barsXYs := make(plotter.XYs, len(data))
	healthXYs := make(plotter.XYs, len(data))
	bandXYs := make(plotter.XYs, len(data))
	towerXYs := make(plotter.XYs, len(data))

	// Helper to safely cast time
	getTime := func(t int64) float64 {
		return float64(t)
	}

	// Collect Unique Towers for Mapping (needed for y-axis levels)
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

		barsXYs[i].X = t
		barsXYs[i].Y = float64(d.Gateway.Signal.FiveG.Bars)

		healthXYs[i].X = t
		healthXYs[i].Y = analysis.CalculateSignalHealth(d.Gateway.Signal.FiveG.RSRP, d.Gateway.Signal.FiveG.SINR)

		bandXYs[i].X = t
		towerXYs[i].X = t

		// Map bands to levels
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

	// Apply Smoothing / Downsampling
	shouldSmoothBars := false
	if len(data) > 1 {
		startTime := data[0].Gateway.Time.LocalTime
		endTime := data[len(data)-1].Gateway.Time.LocalTime
		duration := float64(endTime - startTime)

		// 2 hours = 7200 seconds
		if duration > 7200 && len(barsXYs) > 300 {
			shouldSmoothBars = true
			barsXYs = downsample(barsXYs, 300)
			healthXYs = downsample(healthXYs, 300)
		}

		// 24 hours = 86400 seconds
		if duration > 86400 && len(latencyXYs) > 600 {
			latencyXYs = downsample(latencyXYs, 600)
			stdDevXYs = downsample(stdDevXYs, 600)
			lossXYs = downsample(lossXYs, 600)
			rsrpXYs = downsample(rsrpXYs, 600)
			sinrXYs = downsample(sinrXYs, 600)
			bandXYs = downsample(bandXYs, 600)
			towerXYs = downsample(towerXYs, 600)
		}
	}

	// Calculate Common Time Axis with Padding
	var minTime, maxTime float64
	if len(data) > 0 {
		minTime = getTime(data[0].Gateway.Time.LocalTime)
		maxTime = getTime(data[len(data)-1].Gateway.Time.LocalTime)
	}

	// Add 12% padding to the right for labels
	duration := maxTime - minTime
	if duration <= 0 {
		duration = 60
	}
	maxX := maxTime + (duration * 0.12)
	minX := minTime

	// Common Time Ticker
	timeTicks := plot.TimeTicks{
		Format: "15:04",
		Time: func(t float64) time.Time {
			return time.Unix(int64(t), 0).Local()
		},
	}

	// ---------------------------
	// 2. Chart Generation
	// ---------------------------

	// --- 1. Latency & Packet Loss ---
	pLat := plot.New()
	pLat.Title.Text = "Latency & Packet Loss"
	pLat.Y.Label.Text = "ms / %"
	pLat.X.Tick.Marker = timeTicks
	pLat.X.Min = minX
	pLat.X.Max = maxX
	pLat.Y.Scale = plot.LogScale{}
	pLat.Y.Tick.Marker = logTicks{}

	lineLat, _ := plotter.NewLine(latencyXYs)
	lineLat.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue

	lineStd, _ := plotter.NewLine(stdDevXYs)
	lineStd.Color = color.RGBA{R: 255, G: 140, B: 0, A: 255} // Dark Orange

	// Filter loss points for scatter plot (only show > 0%)
	var lossScatterData plotter.XYs
	for _, pt := range lossXYs {
		// In the sanitization step above, 0.0 was converted to 0.1.
		// Real loss (e.g. 1/20 = 5%) is 5.0.
		// So we only show dots if Y > 0.11 (safe margin).
		if pt.Y > 0.11 {
			lossScatterData = append(lossScatterData, pt)
		}
	}

	scatterLoss, _ := plotter.NewScatter(lossScatterData)
	scatterLoss.GlyphStyle.Shape = draw.CircleGlyph{}
	scatterLoss.GlyphStyle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black
	scatterLoss.GlyphStyle.Radius = vg.Points(2.5)

	pLat.Add(lineLat, lineStd, scatterLoss)
	pLat.Legend.Add("Avg (ms)", lineLat)
	pLat.Legend.Add("StdDev", lineStd)
	pLat.Legend.Add("Loss (%)", scatterLoss)
	pLat.Add(plotter.NewGrid())

	// Add Labels
	addLastPointLabel(pLat, latencyXYs, "%.0fms", lineLat.Color)
	addLastPointLabel(pLat, stdDevXYs, "%.0f", lineStd.Color)
	// Optionally label Loss if > 0.1 (sanitized 0)
	if len(lossXYs) > 0 && lossXYs[len(lossXYs)-1].Y > 0.11 {
		addLastPointLabel(pLat, lossXYs, "%.1f%%", scatterLoss.GlyphStyle.Color)
	}

	// --- 2. Signal Strength (Split) ---
	
	// 2a. SINR (Top)
	pSINR := plot.New()
	pSINR.Title.Text = "Signal Strength (SINR / RSRP)"
	pSINR.Y.Label.Text = "SINR (dB)"
	pSINR.X.Tick.Marker = timeTicks
	pSINR.X.Tick.Label.Color = color.Transparent // Hide labels
	pSINR.X.Min = minX
	pSINR.X.Max = maxX

	lineSINR, _ := plotter.NewLine(sinrXYs)
	lineSINR.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green

	pSINR.Add(lineSINR)
	pSINR.Legend.Add("SINR", lineSINR)
	pSINR.Add(plotter.NewGrid())
	addLastPointLabel(pSINR, sinrXYs, "%.1f", lineSINR.Color)

	// 2b. RSRP (Bottom)
	pRSRP := plot.New()
	pRSRP.Y.Label.Text = "RSRP (dBm)"
	pRSRP.X.Tick.Marker = timeTicks
	pRSRP.X.Min = minX
	pRSRP.X.Max = maxX

	lineRSRP, _ := plotter.NewLine(rsrpXYs)
	lineRSRP.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red

	pRSRP.Add(lineRSRP)
	pRSRP.Legend.Add("RSRP", lineRSRP)
	pRSRP.Add(plotter.NewGrid())
	addLastPointLabel(pRSRP, rsrpXYs, "%.0f", lineRSRP.Color)

	// --- 3. 5G Band & Tower Plot ---
	pBand := plot.New()
	pBand.Title.Text = "Connection Info / Tower & Radio Bands"
	pBand.X.Tick.Marker = timeTicks
	pBand.X.Min = minX
	pBand.X.Max = maxX

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

	lineBand, _ := plotter.NewLine(bandXYs)
	lineBand.StepStyle = plotter.PreStep
	lineBand.Color = color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange

	lineTower, _ := plotter.NewLine(towerXYs)
	lineTower.StepStyle = plotter.PreStep
	lineTower.Color = color.RGBA{R: 255, G: 0, B: 255, A: 255} // Magenta

	pBand.Add(lineBand, lineTower)
	pBand.Add(plotter.NewGrid())

	// --- 4. Signal Bars Plot ---
	pBars := plot.New()
	pBars.Title.Text = "Signal Bars"
	pBars.Y.Label.Text = "Bars"
	pBars.X.Tick.Marker = timeTicks
	pBars.X.Min = minX
	pBars.X.Max = maxX
	pBars.Y.Min = 0
	pBars.Y.Max = 5
	pBars.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: 1, Label: "1"},
		{Value: 2, Label: "2"},
		{Value: 3, Label: "3"},
		{Value: 4, Label: "4"},
	})

	// Health Area (Background)
	healthAreaXYs := make(plotter.XYs, len(healthXYs)+2)
	copy(healthAreaXYs[1:], healthXYs)
	healthAreaXYs[0] = plotter.XY{X: healthXYs[0].X, Y: 0}
	healthAreaXYs[len(healthAreaXYs)-1] = plotter.XY{X: healthXYs[len(healthXYs)-1].X, Y: 0}

	polyHealth, _ := plotter.NewPolygon(healthAreaXYs)
	polyHealth.Color = color.RGBA{R: 240, G: 240, B: 240, A: 255} // Very Light Grey
	polyHealth.LineStyle.Width = 0

	lineHealth, _ := plotter.NewLine(healthXYs)
	lineHealth.Color = color.RGBA{R: 169, G: 169, B: 169, A: 255} // Dark Grey
	lineHealth.Width = vg.Points(1)

	lineBars, _ := plotter.NewLine(barsXYs)
	lineBars.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black
	if shouldSmoothBars {
		lineBars.StepStyle = plotter.NoStep
	} else {
		lineBars.StepStyle = plotter.PreStep
	}
	lineBars.Width = vg.Points(0.8)

	pBars.Add(polyHealth, lineHealth, lineBars)
	pBars.Legend.Add("Reported Bars", lineBars)
	pBars.Legend.Add("Signal Health", lineHealth)
	pBars.Add(plotter.NewGrid())
	addLastPointLabel(pBars, barsXYs, "%.1f", lineBars.Color)


	// ---------------------------
	// 3. Layout & Drawing
	// ---------------------------
	const width = 20 * vg.Inch
	const height = 10 * vg.Inch

	c := vgimg.NewWith(vgimg.UseWH(width, height), vgimg.UseBackgroundColor(color.White))
	dc := draw.New(c)

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
	rectSINR := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: colWidth, Y: rowHeight + (rowHeight / 2)},
			Max: vg.Point{X: width, Y: height},
		},
	}
	pSINR.Draw(rectSINR)

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
