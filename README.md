# Signal Sentry

Signal Sentry is a CLI tool designed to monitor and display real-time signal statistics from a T-Mobile Home Internet gateway. It helps users optimize their gateway placement and monitor connection stability over time.

## Key Features

- **Real-Time Dashboard:** A live console output with color-coded metrics (Excellent/Fair/Poor) for RSRP, SINR, and more.
- **Placement Optimization:** Instant feedback on signal changes to help identify the best spot for your gateway.
- **Detailed Signal Metrics:** View information about 4G/5G bands, tower identification (gNBID/CID), and signal power/quality.
- **Stability Monitoring:** Track signal performance and tower switching in real-time.

## Prerequisites

- **T-Mobile Home Internet Gateway:** Currently tested with gateways having the local API enabled at `http://192.168.12.1`.
- **Go:** Version 1.25.5 or later (for building from source).

## Installation

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/signal-sentry.git
   cd signal-sentry
   ```

2. Build the binary using the Makefile:
   ```bash
   make build
   ```

## Usage

Run the tool using the compiled binary:
```bash
./tmobile-stats
```

Or run it directly using Go:
```bash
make run
```

## Dashboard Preview

The tool provides a live, color-coded dashboard (colors not shown in this text preview) to help you hunt for the best signal.

```text
====================================================================================================
 DEVICE INFO | Model: TMO-G5AR   | FW: 1.00.02    | Serial: XXXXXXXXXXX     | MAC: XX:XX:XX:XX:XX:XX
====================================================================================================

SIGNAL METRICS GUIDE:
---------------------
* BAND:  The frequency band in use (e.g., n41 is faster mid-band, n71 is longer range).
* RSRP:  (Reference Signal Received Power) Your main signal strength.
         Excellent > -80 | Good -80 to -90 | Fair -90 to -100 | Poor < -100.
* SINR:  (Signal-to-Interference-plus-Noise Ratio) Your signal quality (cleanliness).
         Higher is better. > 20 is excellent. < 0 means high noise/interference.
* BARS:  General signal rating (1-5).

 TYPE | BANDS      | BARS | RSRP | SINR | RSRQ | RSSI | CID   | TOWER (gNBID/PCID)
------+------------+------+------+------+------+------+-------+-------------------
 5G   | n41        | 2.0  | -106 | -1   | -14  | -91  | 302   | 1870191  <-- Weak signal in corner
 5G   | n41        | 2.5  | -102 | 2    | -13  | -90  | 302   | 1870191  <-- Moving towards window...
 5G   | n41        | 3.0  | -95  | 8    | -12  | -88  | 302   | 1870191
 5G   | n41        | 4.0  | -88  | 15   | -11  | -85  | 302   | 1870191  <-- Getting better
 5G   | n41        | 5.0  | -79  | 25   | -10  | -82  | 302   | 1870191  <-- Excellent placement found!
```

### Signal Metrics Guide

- **BAND:** The frequency band in use (e.g., n41, n71).
- **RSRP (Reference Signal Received Power):** Main signal strength.
  - Excellent: > -80
  - Good: -80 to -90
  - Fair: -90 to -100
  - Poor: < -100
- **SINR (Signal-to-Interference-plus-Noise Ratio):** Signal quality.
  - Excellent: > 20
  - Poor: < 0
- **CID & gNBID:** Tower and sector identifiers. Changes here indicate tower switching.

## Development

- **Run tests:** `make test`
- **Clean build artifacts:** `make clean`

## License

This project is licensed under the GNU General Public License v2.0 - see the [LICENSE](LICENSE) file for details.
