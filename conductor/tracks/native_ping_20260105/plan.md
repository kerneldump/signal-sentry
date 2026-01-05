# Track Plan: Native Ping Implementation

## Phase 1: Dependency & Setup
- [x] Task: Add `github.com/prometheus-community/pro-bing` dependency.
- [x] Task: Verify: Ensure the library installs correctly.

## Phase 2: Logic Refactoring
- [x] Task: Modify `internal/pinger/pinger.go` to import the new library.
- [x] Task: Rewrite `Run()` and `ping()` to instantiate a `pro-bing` pinger.
    -   Set `Count = 1`.
    -   Set `Timeout = p.Interval`.
    -   On finish, extract RTT and update internal `models.PingStats`.
- [x] Task: Ensure Welford's algorithm (or the library's built-in stats) is used to maintain the cumulative/interval stats correctly.
    -   *Decision:* The library provides stats for the *batch* (which is just 1 ping). We still need our own aggregator for the "Interval" and "Lifetime" windows.

## Phase 3: Verification
- [x] Task: Build the application.
- [x] Task: Run with `sudo ./tmobile-stats -live` and verify ping data is populating.
- [x] Task: Verify packet loss detection (e.g., by pinging a non-existent IP).