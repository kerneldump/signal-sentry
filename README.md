# Signal Sentry

Signal Sentry is a CLI tool designed to monitor and display real-time signal statistics from a T-Mobile Home Internet gateway. It helps users optimize their gateway placement, monitor connection stability, and analyze historical performance.

## Key Features

- **Interactive Live Dashboard:** A TUI (Text User Interface) with real-time graphs, color-coded metrics, and dynamic controls.
- **Native Ping Integration:** Monitors latency and packet loss alongside signal stats (requires `sudo`).
- **Historical Analysis:** Analyze collected logs to generate summaries of signal quality, bands, and tower usage over time.
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
sudo ./tmobile-stats [flags] [subcommand]
```

### Common Workflows

**1. Live Monitoring (Recommended)**
Start the interactive TUI dashboard. This will also log data to `stats.log` in the background.
```bash
sudo ./tmobile-stats -live
```
*   **Controls:**
    *   `+` / `-`: Increase/Decrease refresh interval.
    *   `i`: Toggle the help overlay.
    *   `q`: Quit.

**2. Historical Analysis**
Generate a report from your collected logs.
```bash
./tmobile-stats analyze
```
*   *Note: You can run this in a separate terminal while the monitor is running.*

**3. Legacy/Scripting Mode**
Run with standard standard output (useful for piping to other tools).
```bash
sudo ./tmobile-stats -interval 2
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

## Configuration

A `config.json` file can be used to set defaults. Example:

```json
{
  "router_url": "http://192.168.12.1/TMI/v1/gateway?get=all",
  "ping_target": "8.8.8.8",
  "refresh_interval": 5,
  "live_mode": true
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
         n25: Balanced speed and range.
         n71: Long range, slower speeds.

* RSRP:  (Reference Signal Received Power) Your main signal strength.
         Excellent > -80  | Good -80 to -95
         Fair -95 to -110 | Poor < -110 (Risk of drops).

* SINR:  (Signal-to-Interference-plus-Noise Ratio) Signal quality.
         Higher is better. > 20 is excellent.
         < 0 means high noise (your speed will suffer).

 TYPE  | BANDS      | BARS | RSRP  | SINR  | RSRQ | RSSI | CID   | TWR gNBID/PCIDE | MIN AVG MAX STD LOSS
-------+------------+------+-------+-------+------+------+-------+-----------------+-------------------------
  5G   | n41        | 4.0  | -88   | 15    | -11  | -85  | 302   | 1870191         | 18.2 22.5 45.1 5.2 0.0%
  5G   | n41        | 5.0  | -79   | 25    | -10  | -82  | 302   | 1870191         | 17.8 21.0 28.4 2.1 0.0%
```

## Development

- **Run tests:** `make test`
- **Clean build artifacts:** `make clean`

## License

This project is licensed under the GNU General Public License v2.0 - see the [LICENSE](LICENSE) file for details.
