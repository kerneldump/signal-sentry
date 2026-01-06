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
	pLat := plot.New()
	pLat.Title.Text = "Latency & StdDev"
	pLat.Y.Label.Text = "ms"
	pLat.X.Tick.Marker = timeTicks

	lineLat, _ := plotter.NewLine(latencyXYs)
	lineLat.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue

	lineStd, _ := plotter.NewLine(stdDevXYs)
	lineStd.Color = color.RGBA{R: 0, G: 255, B: 255, A: 255}           // Cyan
	lineStd.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)} // Dashed

	pLat.Add(lineLat, lineStd)
	pLat.Legend.Add("Avg", lineLat)
	pLat.Legend.Add("StdDev", lineStd)
	pLat.Add(plotter.NewGrid())

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

	pLoss := plot.New()
	pLoss.Title.Text = "Packet Loss"
	pLoss.Y.Label.Text = "%"
	pLoss.X.Tick.Marker = timeTicks
	pLoss.Y.Min = 0
	pLoss.Y.Max = 100

	lineLoss, _ := plotter.NewLine(lossXYs)
	lineLoss.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black

	pLoss.Add(lineLoss)
	pLoss.Add(plotter.NewGrid())

	// Combine into a single image
	const width = 10 * vg.Inch
	const height = 12 * vg.Inch

	c := vgimg.NewWith(vgimg.UseWH(width, height), vgimg.UseBackgroundColor(color.White))
	dc := draw.New(c)

	// Layout: 3 rows, 1 column

	rectLat := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: height * 2 / 3},
			Max: vg.Point{X: width, Y: height},
		},
	}
	pLat.Draw(rectLat)

	rectSig := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: height * 1 / 3},
			Max: vg.Point{X: width, Y: height * 2 / 3},
		},
	}
	pH.Draw(rectSig)

	rectLoss := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: width, Y: height * 1 / 3},
		},
	}
	pLoss.Draw(rectLoss)

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
