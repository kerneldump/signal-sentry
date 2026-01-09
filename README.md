# Signal Sentry

Signal Sentry is a CLI tool designed to monitor and display real-time signal statistics from a T-Mobile Home Internet gateway. It helps users optimize their gateway placement, monitor connection stability, and analyze historical performance.

## Key Features

- **Interactive Live Dashboard:** A TUI (Text User Interface) with real-time graphs, color-coded metrics, and dynamic controls.
- **Live Web Dashboard:** A lightweight local web server to view charts in your browser with selectable time ranges and auto-refresh.
- **Native Ping Integration:** Monitors latency and packet loss alongside signal stats (requires `sudo`).
- **Smart Historical Analysis:** Generate detailed reports with time-based filtering (`-range 24h`) and "Signal Health" scoring.
- **Advanced Charting:** Creates high-resolution 2x2 grid charts visualizing Signal Strength, Latency, Bands, and Signal Bars vs Health. Automatically smooths data for long-term trends.
- **Placement Optimization:** Instant feedback on signal changes to help identify the best spot for your gateway.
- **Detailed Signal Metrics:** View information about 4G/5G bands, tower identification (gNBID/CID), RSRP, SINR, and more.
- **Automatic Logging:** Silently records all data to `stats.log` for future analysis.

## Prerequisites

- **T-Mobile Home Internet Gateway:** Currently tested with gateways having the local API enabled at `http://192.168.12.1`.
- **Root Privileges:** Required for ICMP ping functionality (use `sudo`).
- **Go:** Version 1.25.5 or later (for building from source).

## Installation

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/signal-sentry.git
   cd signal-sentry
   ```

2. Build the binary:
   ```bash
   make build
   ```

## Usage

Run the tool using the compiled binary. **Note: `sudo` is required for ping statistics.**

```bash
sudo ./signal-sentry [flags] [subcommand]
```

### Common Workflows

**1. Live Monitoring (Recommended)**
Start the interactive TUI dashboard. This will also log data to `stats.log` in the background.
```bash
sudo ./signal-sentry -live
```
*   **Controls:**
    *   `+` / `-`: Increase/Decrease refresh interval.
    *   `i`: Toggle the help overlay.
    *   `q`: Quit.

**2. Historical Analysis**
Generate a report from your collected logs.
```bash
# Analyze all data
./signal-sentry analyze

# Analyze last 24 hours
./signal-sentry analyze -range 24h
```
*   *Note: You can run this in a separate terminal while the monitor is running.*

**3. Generate Charts**
Create a visual graph of your signal history.
```bash
# Generate smoothed chart for long duration
./signal-sentry chart -input stats.log -output my-signal.png

# Generate detailed chart for specific incident
./signal-sentry chart -start "2026-01-05 14:00:00" -end "2026-01-05 16:00:00"
```

**4. Live Web Dashboard**
View your signal charts in any web browser. This starts a local server that serves auto-refreshing charts.
```bash
./signal-sentry web
```
*   **Access:** Open `http://localhost:8080` in your browser.
*   **Features:** Toggle between 1h, 6h, 24h, or Max history views instantly.

**5. Legacy/Scripting Mode**
Run with standard standard output (useful for piping to other tools).
```bash
sudo ./signal-sentry -interval 2
```

### Flags

- `-live`: Enable the interactive TUI dashboard.
- `-interval int`: Refresh interval in seconds (default: 5).
- `-config string`: Path to config file (default: `config.json`).
- `-no-auto-log`: Disable automatic logging to `stats.log` (useful if you are running a second instance just to view).
- `-format string`: Output format for *additional* file logging (`json` or `csv`).
- `-output string`: Output filename for the formatted log.
- `-version`: Show version information.

### Subcommands

- `analyze`: Parse a log file and display summary statistics.
  - `-input`: Path to the log file (default: `stats.log`).
  - `-range`: Relative time range from now (e.g., `24h`, `30m`).
  - `-start`: Start date/time (format: `YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`).
  - `-end`: End date/time (format: `YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`).
- `chart`: Generate a PNG chart of RSRP and SINR over time from a log file.
  - `-input`: Path to the log file (default: `stats.log`).
  - `-output`: Path to save the chart image (default: `signal-analysis.png`).
  - `-range`, `-start`, `-end`: Same filtering options as `analyze`.
- `web`: Start a local web server to view auto-refreshing signal charts.
  - `-port`: Port to listen on (default: `8080`).
  - `-input`: Path to the log file (default: `stats.log`).

### Configuration

A `config.json` file can be used to set defaults. Example:

```json
{
  "router_url": "http://192.168.12.1/TMI/v1/gateway?get=all",
  "ping_target": "8.8.8.8",
  "refresh_interval": 5,
  "live_mode": true,
  "format": "csv",
  "output": "my-log.csv",
  "disable_auto_log": false
}
```

## Dashboard Preview

The tool provides a live, color-coded dashboard to help you hunt for the best signal.

```text
====================================================================================================
 DEVICE INFO | Model: TMO-G5AR   | FW: 1.00.02    | Serial: XXXXXXXXXXX     | MAC: XX:XX:XX:XX:XX:XX
====================================================================================================

SIGNAL METRICS GUIDE:
---------------------
* BAND:  The frequency band in use.
         n41: High speed, shorter range (Ultra Capacity).
...
```

### Analysis Report Example

```text
================================================================================
 HISTORICAL SIGNAL ANALYSIS
================================================================================
Filter:        2026-01-06 15:00:39 to 2026-01-06 16:00:39
Data Range:    2026-01-06 15:00:41 to 16:00:31
Duration:      59m50s
Total Samples: 360

METRIC      MIN   AVG    MAX
------      ---   ---    ---
RSRP (dBm)  -104  -98.5  -95
SINR (dB)   2     7.7    16
Ping (ms)   18.0  32.5   561.0

RELIABILITY:
  Packet Loss: 0 / 360 (0.00%)

BANDS SEEN:
  n41        360 samples (100.0%)    59m 50s

TOWERS SEEN:
  1870191    360 samples (100.0%) live    59m 50s

BARS SEEN:
  3          79 samples (21.9%)    13m 6s
  4          281 samples (78.1%) real-time    46m 44s

BARS AVG:
Overall     3.8
Last 1h     3.8
SgnlHealth  3.7
================================================================================
```

## Development

- **Run tests:** `make test`
- **Clean build artifacts:** `make clean`

## License

This project is licensed under the GNU General Public License v2.0 - see the [LICENSE](LICENSE) file for details.
