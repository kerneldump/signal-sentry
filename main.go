package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"tmobile-stats/internal/config"
	"tmobile-stats/internal/gateway"
	"tmobile-stats/internal/logger"
	"tmobile-stats/internal/models"
	"tmobile-stats/internal/pinger"
)

const (
	headerInterval = 20
	Version        = "v1.1.0"
)

func main() {
	// 1. Initial Flags for Config Loading
	configPath := flag.String("config", "", "Path to config file (JSON)")
	
	// Temporarily define other flags to avoid parsing errors
	intervalFlag := flag.Int("interval", 0, "Refresh interval in seconds")
	formatFlag := flag.String("format", "", "Output format (json or csv)")
	outputFlag := flag.String("output", "", "Output filename")
	versionFlag := flag.Bool("version", false, "Show version information")
	liveFlag := flag.Bool("live", false, "Enable interactive live view")
	noAutoLogFlag := flag.Bool("no-auto-log", false, "Disable automatic logging to stats.log")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Signal Sentry - T-Mobile Gateway Signal Monitor (%s)\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage:\n  signal-sentry [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Signal Sentry %s\n", Version)
		os.Exit(0)
	}

	// 2. Load Config
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 3. Override Config with explicitly provided flags
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "interval":
			cfg.RefreshInterval = *intervalFlag
		case "format":
			cfg.Format = *formatFlag
		case "output":
			cfg.Output = *outputFlag
		case "live":
			cfg.LiveMode = *liveFlag
		case "no-auto-log":
			cfg.DisableAutoLog = *noAutoLogFlag
		}
	})

	// 4. Validate Final Config
	if err := validateInterval(cfg.RefreshInterval); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := validateFormat(cfg.Format); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Determine default filename if not provided but format is set
	if cfg.Format != "" && cfg.Output == "" {
		if cfg.Format == "json" {
			cfg.Output = "signal-data.json"
		} else {
			cfg.Output = "signal-data.csv"
		}
	}

	// Initialize User Logger
	var appLogger logger.Logger
	if cfg.Format != "" {
		var err error
		if cfg.Format == "json" {
			appLogger, err = logger.NewJSONLogger(cfg.Output)
		} else if cfg.Format == "csv" {
			appLogger, err = logger.NewCSVLogger(cfg.Output)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer appLogger.Close()
	}

	// Initialize Background Logger (Always-on JSON)
	var bgLogger logger.Logger
	if !cfg.DisableAutoLog {
		// Avoid double-logging if user explicitly chose stats.log
		if cfg.Output != "stats.log" {
			var err error
			bgLogger, err = logger.NewJSONLogger("stats.log")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize background logger: %v\n", err)
				// We don't exit here, just warn and continue without BG logging
			} else {
				defer bgLogger.Close()
			}
		}
	}

	// 5. Initialize Pinger
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pg := pinger.NewPinger(cfg.PingTarget, 1*time.Second)
	go pg.Run(ctx)

	refreshDuration := time.Duration(cfg.RefreshInterval) * time.Second

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	firstRun := true
	linesPrinted := 0

	for {
		// 1. Fetch Data
		gatewayData, err := gateway.FetchStats(client, cfg.RouterURL)
		if err != nil {
			fmt.Printf("Error fetching stats: %v\n", err)
			time.Sleep(refreshDuration)
			continue
		}

		// Get stats for the specific interval window (and reset for next window)
		pingData := pg.GetStatsAndReset()
		
		data := &models.CombinedStats{
			Gateway: *gatewayData,
			Ping:    pingData,
		}

		// 2. Log Data if enabled (User + Background)
		if appLogger != nil {
			if err := appLogger.Log(data); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to log data: %v\n", err)
			}
		}
		if bgLogger != nil {
			if err := bgLogger.Log(data); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to stats.log: %v\n", err)
			}
		}

		// 3. Initial Setup (Device Info & Legend)
		if firstRun {
			printDeviceInfo(data.Gateway.Device)
			printLegend()
			firstRun = false
		}

		// 3. Print Header periodically
		if linesPrinted%headerInterval == 0 {
			printHeader()
		}

		// 4. Print Data Rows
		printRow("5G", data.Gateway.Signal.FiveG, data.Ping)
		// Check if we have valid 4G data (e.g., non-zero bars or bands)
		if len(data.Gateway.Signal.FourG.Bands) > 0 || data.Gateway.Signal.FourG.Bars > 0 {
			printRow("4G", data.Gateway.Signal.FourG, data.Ping)
			linesPrinted++ // Count extra line for 4G
		}

		linesPrinted++
		time.Sleep(refreshDuration)
	}
}

func printDeviceInfo(d models.DeviceInfo) {
	fmt.Println("================================================================================================================================================================")
	fmt.Printf(" DEVICE INFO | Model: %-10s | FW: %-10s | Serial: %-15s | MAC: %s\n",
		d.Model, d.SoftwareVersion, d.Serial, d.MacID)
	fmt.Println("================================================================================================================================================================")
}

func printLegend() {
	fmt.Print(`
SIGNAL METRICS GUIDE:
---------------------
* BAND:  The frequency band in use.
         n41: High speed, shorter range (Ultra Capacity).
         n25: Balanced speed and range.
         n71: Long range, slower speeds.

* RSRP:  (Reference Signal Received Power) Your main signal strength.
         Excellent > -80  | Good -80 to -95
         Fair -95 to -110 | Poor < -110 (Risk of drops).

* SINR:  (Signal-to-Interference-plus-Noise Ratio) Signal quality.
         Higher is better. > 20 is excellent.
         < 0 means high noise (your speed will suffer).

* RSRQ:  (Reference Signal Received Quality) The congestion indicator.
         If SINR is Good (high) but RSRQ is Bad (low), the tower is
         likely congested with heavy traffic.

* CID & gNBID:
         gNBID identifies the physical TOWER.
         CID identifies the SECTOR (which side of the tower you are facing).
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

func printHeader() {
	// Column Widths (Visible):
	// Type: 6, Bands: 12, Bars: 6, RSRP: 7, SINR: 7, RSRQ: 6, RSSI: 6, CID: 7, Tower: 17, PingLoss: 25
	fmt.Println(" TYPE  | BANDS      | BARS | RSRP  | SINR  | RSRQ | RSSI | CID   | TWR gNBID/PCIDE | PING MIN AVG MAX STD LOSS")
	fmt.Println("-------+------------+------+-------+-------+------+------+-------+-----------------+-------------------------")
}

func printRow(connType string, stats models.ConnectionStats, ping models.PingStats) {
	bands := strings.Join(stats.Bands, ",")
	if bands == "" {
		bands = "---"
	}

	towerID := stats.GNBID
	if towerID == 0 {
		towerID = stats.PCID
	}

	// Colorize Connection Type
	typeStr := connType
	if connType == "5G" {
		typeStr = fmt.Sprintf("%s%-2s%s", ColorCyan, connType, ColorReset)
	} else if connType == "4G" {
		typeStr = fmt.Sprintf("%s%-2s%s", ColorBlue, connType, ColorReset)
	}

	// Format each value to exactly 4 characters (e.g. 44.0 or 0.0%)
	// Note: Loss formatted as %3.1f to be e.g. "0.0" then we append "%"
	pMin := fmt.Sprintf("%4.1f", ping.Min)
	pAvg := fmt.Sprintf("%4.1f", ping.Avg)
	pMax := fmt.Sprintf("%4.1f", ping.Max)
	pStd := fmt.Sprintf("%4.1f", ping.StdDev)
	
	lossVal := fmt.Sprintf("%3.1f%%", ping.Loss)
	if ping.Loss > 0 {
		lossVal = fmt.Sprintf("%s%s%s", ColorRed, lossVal, ColorReset)
	}

	// Print row with aligned columns
	// Combined Ping/Loss section starts after one space following the pipe
	fmt.Printf("  %s   | %-10s | %s  | %s  | %s  | %-4d | %-4d | %-5d | %-15d | %s %s %s %s %s \n",
		typeStr,
		bands,
		colorizeBars(stats.Bars),
		colorizeRSRP(stats.RSRP),
		colorizeSINR(stats.SINR),
		stats.RSRQ,
		stats.RSSI,
		stats.CID,
		towerID,
		pMin, pAvg, pMax, pStd, lossVal,
	)
}

func colorizeRSRP(val int) string {
	s := fmt.Sprintf("%4d", val)
	if val > -80 {
		return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorReset)
	} else if val >= -110 {
		return fmt.Sprintf("%s%s%s", ColorYellow, s, ColorReset)
	}
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorReset)
}

func colorizeSINR(val int) string {
	s := fmt.Sprintf("%4d", val)
	if val > 20 {
		return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorReset)
	} else if val >= 0 {
		return fmt.Sprintf("%s%s%s", ColorYellow, s, ColorReset)
	}
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorReset)
}

func colorizeBars(val float64) string {
	s := fmt.Sprintf("%3.1f", val)
	if val >= 4.0 {
		return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorReset)
	} else if val >= 2.0 {
		return fmt.Sprintf("%s%s%s", ColorYellow, s, ColorReset)
	}
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorReset)
}