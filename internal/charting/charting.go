package charting

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"time"

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
	rsrpXYs := make(plotter.XYs, len(data))
	sinrXYs := make(plotter.XYs, len(data))
	lossXYs := make(plotter.XYs, len(data))

	startTime := time.Unix(data[0].Gateway.Time.LocalTime, 0)

	for i, d := range data {
		t := time.Unix(d.Gateway.Time.LocalTime, 0).Sub(startTime).Minutes()

		latencyXYs[i].X = t
		latencyXYs[i].Y = d.Ping.Avg // Using Avg for the main line

		rsrpXYs[i].X = t
		rsrpXYs[i].Y = float64(d.Gateway.Signal.FiveG.RSRP)

		sinrXYs[i].X = t
		sinrXYs[i].Y = float64(d.Gateway.Signal.FiveG.SINR)

		lossXYs[i].X = t
		lossXYs[i].Y = d.Ping.Loss
	}

	// Create Plots
	pLat := plot.New()
	pLat.Title.Text = "Latency (AVG)"
	pLat.Y.Label.Text = "ms"
	pLat.X.Label.Text = "Time (min)"
	lineLat, _ := plotter.NewLine(latencyXYs)
	lineLat.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue
	pLat.Add(lineLat)
	pLat.Add(plotter.NewGrid())

	pH := plot.New()
	pH.Title.Text = "Signal Strength"
	pH.Y.Label.Text = "dBm / dB"
	pH.X.Label.Text = "Time (min)"
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
	pLoss.X.Label.Text = "Time (min)"
	pLoss.Y.Min = 0
	pLoss.Y.Max = 100
	lineLoss, _ := plotter.NewLine(lossXYs)
	lineLoss.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black

	// Create a filled polygon for loss to make it stand out
	// We'd need to construct the full polygon loop, but line is fine for now for MVP.

	pLoss.Add(lineLoss)
	pLoss.Add(plotter.NewGrid())

	// Combine into a single image using a Tiles arrangement or just create 3 separate plots drawn to one canvas
	// Gonum doesn't have a built-in "subplot" layout system as simple as Matplotlib,
	// so we create a large canvas and draw each plot onto sub-tiles.

	const width = 10 * vg.Inch
	const height = 12 * vg.Inch

	// Create a canvas with white background
	c := vgimg.NewWith(vgimg.UseWH(width, height), vgimg.UseBackgroundColor(color.White))
	dc := draw.New(c)

	// Layout: 3 rows, 1 column
	// Top: Latency
	// Mid: Signal
	// Bot: Loss

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
