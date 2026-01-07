# Track Plan: Improve Logarithmic Axis Labels

## Phase 1: Implementation
- [x] Task: Implement custom tick marker in `internal/charting/charting.go` to format log labels as integers.
    -   Define a `logTicks` struct that implements `plot.Ticker`.
    -   Use `plot.LogTicks{}` internally to get major ticks.
    -   Format labels using `fmt.Sprintf("%.0f", val)`.
- [x] Task: Apply the custom ticker to `pLat.Y.Tick.Marker`.

## Phase 2: Verification
- [x] Task: Verify: Generate a chart and ensure labels are readable integers.
