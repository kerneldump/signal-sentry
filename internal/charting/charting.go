package charting

import (
	"fmt"
	"image/color"
	"image/png"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"tmobile-stats/internal/analysis"
	"tmobile-stats/internal/models"
)

// Generate creates a PNG chart from the provided stats and saves it to outputFile.
func Generate(data []models.CombinedStats, outputFile string) error {
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

		latencyXYs[i].X = t
		latencyXYs[i].Y = d.Ping.Avg

		stdDevXYs[i].X = t
		stdDevXYs[i].Y = d.Ping.StdDev

		rsrpXYs[i].X = t
		rsrpXYs[i].Y = float64(d.Gateway.Signal.FiveG.RSRP)

		sinrXYs[i].X = t
		sinrXYs[i].Y = float64(d.Gateway.Signal.FiveG.SINR)

		lossXYs[i].X = t
		lossXYs[i].Y = d.Ping.Loss
	}

	// Common Time Ticker
	// We use the standard TimeTicks which interprets float64 X values as unix seconds
	timeTicks := plot.TimeTicks{Format: "15:04"}

	// Create Plots

	// 1. Latency & Packet Loss
	pLat := plot.New()
	pLat.Title.Text = "Latency & Packet Loss"
	pLat.Y.Label.Text = "ms / %"
	pLat.X.Tick.Marker = timeTicks

	lineLat, _ := plotter.NewLine(latencyXYs)
	lineLat.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue

	lineStd, _ := plotter.NewLine(stdDevXYs)
	lineStd.Color = color.RGBA{R: 0, G: 255, B: 255, A: 255}           // Cyan
	lineStd.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)} // Dashed

	lineLoss, _ := plotter.NewLine(lossXYs)
	lineLoss.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black

	pLat.Add(lineLat, lineStd, lineLoss)
	pLat.Legend.Add("Avg (ms)", lineLat)
	pLat.Legend.Add("StdDev", lineStd)
	pLat.Legend.Add("Loss (%)", lineLoss)
	pLat.Add(plotter.NewGrid())

	// 2. Signal Strength
	pH := plot.New()
	pH.Title.Text = "Signal Strength"
	pH.Y.Label.Text = "dBm / dB"
	pH.X.Tick.Marker = timeTicks

	lineRSRP, _ := plotter.NewLine(rsrpXYs)
	lineRSRP.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red

	lineSINR, _ := plotter.NewLine(sinrXYs)
	lineSINR.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green

	pH.Add(lineRSRP, lineSINR)
	pH.Legend.Add("RSRP", lineRSRP)
	pH.Legend.Add("SINR", lineSINR)
	pH.Add(plotter.NewGrid())

	// 3. 5G Band Plot
	pBand := plot.New()
	pBand.Title.Text = "5G Band"
	pBand.X.Tick.Marker = timeTicks

	// Custom Y Ticks for Bands
	pBand.Y.Min = 0
	pBand.Y.Max = 4
	pBand.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: 1, Label: "n71 (Range)"},
		{Value: 2, Label: "n25 (Mid)"},
		{Value: 3, Label: "n41 (Speed)"},
	})

	bandXYs := make(plotter.XYs, len(data))
	for i, d := range data {
		t := getTime(d.Gateway.Time.LocalTime)
		bandXYs[i].X = t

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
	}

	lineBand, _ := plotter.NewLine(bandXYs)
	lineBand.StepStyle = plotter.PreStep
	lineBand.Color = color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange

	pBand.Add(lineBand)
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

	lineBars, _ := plotter.NewLine(barsXYs)
	lineBars.StepStyle = plotter.PreStep
	lineBars.Color = color.RGBA{R: 128, G: 0, B: 128, A: 255} // Purple

	lineHealth, _ := plotter.NewLine(healthXYs)
	lineHealth.Color = color.RGBA{R: 0, G: 128, B: 128, A: 255} // Teal

	pBars.Add(lineBars, lineHealth)
	pBars.Legend.Add("Reported Bars", lineBars)
	pBars.Legend.Add("Signal Health", lineHealth)
	pBars.Add(plotter.NewGrid())

	// Combine into a single image
	const width = 10 * vg.Inch
	const height = 16 * vg.Inch // Maintains 4 rows height

	c := vgimg.NewWith(vgimg.UseWH(width, height), vgimg.UseBackgroundColor(color.White))
	dc := draw.New(c)

	// Layout: 4 rows, 1 column
	// Row 1 (Top): Latency + Loss
	// Row 2: Signal
	// Row 3: Bars
	// Row 4 (Bot): Bands

	rowHeight := height / 4

	rectLat := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: rowHeight * 3},
			Max: vg.Point{X: width, Y: height},
		},
	}
	pLat.Draw(rectLat)

	rectSig := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: rowHeight * 2},
			Max: vg.Point{X: width, Y: rowHeight * 3},
		},
	}
	pH.Draw(rectSig)

	rectBars := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: rowHeight * 1},
			Max: vg.Point{X: width, Y: rowHeight * 2},
		},
	}
	pBars.Draw(rectBars)

	rectBand := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: width, Y: rowHeight},
		},
	}
	pBand.Draw(rectBand)

	// Save
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, c.Image()); err != nil {
		return err
	}
	return nil
}
