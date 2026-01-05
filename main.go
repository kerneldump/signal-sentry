package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"tmobile-stats/internal/gateway"
	"tmobile-stats/internal/logger"
)

const (
	gatewayURL     = "http://192.168.12.1/TMI/v1/gateway?get=all"
	headerInterval = 20
	Version        = "v1.1.0"
)

func main() {
	intervalPtr := flag.Int("interval", 5, "Refresh interval in seconds")
	formatPtr := flag.String("format", "", "Output format (json or csv)")
	outputPtr := flag.String("output", "", "Output filename (default: signal-data.json/csv)")
	versionPtr := flag.Bool("version", false, "Show version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Signal Sentry - T-Mobile Gateway Signal Monitor (%s)\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage:\n  signal-sentry [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionPtr {
		fmt.Printf("Signal Sentry %s\n", Version)
		os.Exit(0)
	}

	if err := validateInterval(*intervalPtr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := validateFormat(*formatPtr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Determine default filename if not provided but format is set
	var outputFilename string
	if *formatPtr != "" {
		if *outputPtr == "" {
			if *formatPtr == "json" {
				outputFilename = "signal-data.json"
			} else {
				outputFilename = "signal-data.csv"
			}
		} else {
			outputFilename = *outputPtr
		}
	}

	var appLogger logger.Logger
	if *formatPtr != "" {
		var err error
		if *formatPtr == "json" {
			appLogger, err = logger.NewJSONLogger(outputFilename)
		} else if *formatPtr == "csv" {
			appLogger, err = logger.NewCSVLogger(outputFilename)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer appLogger.Close()
	}

	refreshDuration := time.Duration(*intervalPtr) * time.Second

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	firstRun := true
	linesPrinted := 0

	for {
		// 1. Fetch Data
		data, err := gateway.FetchStats(client, gatewayURL)
		if err != nil {
			fmt.Printf("Error fetching stats: %v\n", err)
			time.Sleep(refreshDuration)
			continue
		}

		// 2. Log Data if enabled
		if appLogger != nil {
			if err := appLogger.Log(data); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to log data: %v\n", err)
			}
		}

		// 3. Initial Setup (Device Info & Legend)
		if firstRun {
			printDeviceInfo(data.Device)
			printLegend()
			firstRun = false
		}

		// 3. Print Header periodically
		if linesPrinted%headerInterval == 0 {
			printHeader()
		}

		// 4. Print Data Rows
		printRow("5G", data.Signal.FiveG)
		// Check if we have valid 4G data (e.g., non-zero bars or bands)
		if len(data.Signal.FourG.Bands) > 0 || data.Signal.FourG.Bars > 0 {
			printRow("4G", data.Signal.FourG)
			linesPrinted++ // Count extra line for 4G
		}

		linesPrinted++
		time.Sleep(refreshDuration)
	}
}

func printDeviceInfo(d gateway.DeviceInfo) {
	fmt.Println("====================================================================================================")
	fmt.Printf(" DEVICE INFO | Model: %-10s | FW: %-10s | Serial: %-15s | MAC: %s\n",
		d.Model, d.SoftwareVersion, d.Serial, d.MacID)
	fmt.Println("====================================================================================================")
}

func printLegend() {
	fmt.Print(`
SIGNAL METRICS GUIDE:
---------------------
* BAND:  The frequency band in use (e.g., n41 is faster mid-band, n71 is longer range).
* RSRP:  (Reference Signal Received Power) Your main signal strength.
         Excellent > -80 | Good -80 to -90 | Fair -90 to -100 | Poor < -100.
* SINR:  (Signal-to-Interference-plus-Noise Ratio) Your signal quality (cleanliness).
         Higher is better. > 20 is excellent. < 0 means high noise/interference.
* BARS:  General signal rating (1-5).

* CID & gNBID: 
         These are very important if you are aiming antennas. They tell you which tower and which
         sector of the tower you are talking to. If your signal drops, seeing these numbers change 
         tells you if you switched towers.

* RSRQ & RSSI:
         * RSRQ: Helpful. If SINR is good but RSRQ is bad, the tower is likely congested.
         * RSSI: Less critical for 5G (RSRP is better), but harmless to include.
`)
}

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func colorizeRSRP(val int) string {
	if val > -80 {
		return fmt.Sprintf("%s%d%s", ColorGreen, val, ColorReset)
	} else if val >= -100 {
		return fmt.Sprintf("%s%d%s", ColorYellow, val, ColorReset)
	}
	return fmt.Sprintf("%s%d%s", ColorRed, val, ColorReset)
}

func colorizeSINR(val int) string {
	if val > 20 {
		return fmt.Sprintf("%s%d%s", ColorGreen, val, ColorReset)
	} else if val >= 0 {
		return fmt.Sprintf("%s%d%s", ColorYellow, val, ColorReset)
	}
	return fmt.Sprintf("%s%d%s", ColorRed, val, ColorReset)
}

func colorizeBars(val float64) string {
	if val >= 4.0 {
		return fmt.Sprintf("%s%.1f%s", ColorGreen, val, ColorReset)
	} else if val >= 2.0 {
		return fmt.Sprintf("%s%.1f%s", ColorYellow, val, ColorReset)
	}
	return fmt.Sprintf("%s%.1f%s", ColorRed, val, ColorReset)
}

func printHeader() {
	fmt.Println(" TYPE | BANDS      | BARS | RSRP | SINR | RSRQ | RSSI | CID   | TOWER (gNBID/PCID)")
	fmt.Println("------+------------+------+------+------+------+------+-------+-------------------")
}

func printRow(connType string, stats gateway.ConnectionStats) {
	bands := strings.Join(stats.Bands, ",")
	if bands == "" {
		bands = "---"
	}

	// T-Mobile 5G gateways often return large integers for ID.
	// For 4G, 'PCID' is often the cell ID.
	// We'll prioritize gNBID if non-zero, else PCID, else CID.
	towerID := stats.GNBID
	if towerID == 0 {
		towerID = stats.PCID
	}
	// Note: Sometimes the raw CID is what users look at for cell identity too.

	// Colorize Connection Type
	typeStr := connType
	if connType == "5G" {
		typeStr = fmt.Sprintf("%s%s%s", ColorCyan, connType, ColorReset)
	} else if connType == "4G" {
		typeStr = fmt.Sprintf("%s%s%s", ColorBlue, connType, ColorReset)
	}

	// Adjusted Printf to handle color strings directly (padding won't work perfectly on the string including codes)
	// We will format the NUMBER first, then wrap color.

	fmt.Printf(" %-13s | %-10s | %-13s | %-13s | %-13s | %-4d | %-4d | %-5d | %d\n",
		typeStr,
		bands,
		colorizeBars(stats.Bars),
		colorizeRSRP(stats.RSRP),
		colorizeSINR(stats.SINR),
		stats.RSRQ,
		stats.RSSI,
		stats.CID,
		towerID,
	)
}