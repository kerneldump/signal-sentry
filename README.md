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
