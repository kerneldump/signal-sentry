# Track Plan: Logarithmic Scale for Latency Chart

## Phase 1: Implementation
- [x] Task: Update `Generate` in `internal/charting/charting.go`.
    -   Sanitize data: Ensure Latency, StdDev, and Loss values are >= 0.1 (for Log scale safety).
    -   Configure `pLat.Y.Scale = plot.LogScale{}`.
    -   Configure `pLat.Y.Tick.Marker = plot.LogTicks{}`.

## Phase 2: Verification
- [x] Task: Verify: Generate chart with existing data (`stats.log` has some variation).
- [x] Task: Verify: Check if 0% loss is visible (bottom of log chart).
