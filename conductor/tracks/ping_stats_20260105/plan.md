# Track Plan: Ping Statistics & Advanced Configuration

## Phase 1: Configuration Engine
- [x] Task: Create `internal/config` package and define `Config` struct. [5e3d2a1]
- [x] Task: Implement config file loading (JSON) and default values. [5e3d2a1]
- [x] Task: Refactor `main.go` to use `internal/config` instead of raw flags/constants for URL and Interval. [7a1b2c3]
- [x] Task: Verify: Test loading config from file and overriding with flags. [d4e5f6g]

## Phase 2: Ping Service Implementation
- [x] Task: Create `internal/pinger` package. [h8i9j0k]
- [x] Task: Implement `Run(target string, interval time.Duration, ch chan<- PingStats)` function. [h8i9j0k]
- [x] Task: Implement shell execution of `ping -c 1` and regex parsing for macOS output. [h8i9j0k]
- [x] Task: Calculate/Extract Min, Avg, Max, StdDev, Loss. [h8i9j0k]
- [x] Task: Verify: Write a standalone test/script to verify parsing correctness on the host system. [l1m2n3o]

## Phase 3: Data Aggregation
- [x] Task: Define unified data model (e.g., in `internal/gateway` or `internal/models`) combining Signal + Ping. [p4q5r6s]
- [x] Task: Update `main.go` (and `ui` from Track 1) to listen to both the Signal ticker and the Ping channel (or shared state). [t7u8v9w]
- [x] Task: Ensure thread safety (use mutex if sharing state between routines). [t7u8v9w]

## Phase 4: Output Integration
- [x] Task: Update `internal/logger/csv.go` to add headers and values for Ping stats. [x1y2z3a]
- [x] Task: Update `internal/logger/json.go` to include Ping struct. [x1y2z3a]
- [x] Task: Update CLI output (and TUI View) to display the 5 new columns. [b4c5d6e]
- [x] Task: Verify: Run full app and check all outputs (Screen, CSV, JSON) for data integrity. [f7g8h9i]

