# Track Plan: Bar Averages and Derived Signal Metrics

## Phase 1: Logic Implementation
- [x] Task: Create `internal/analysis/scoring.go`.
    -   Implement `CalculateSignalHealth(rsrp, sinr int) float64` using the 70/30 weighted formula.
    -   Implement unit tests for `CalculateSignalHealth` covering various boundary conditions.
- [x] Task: Update `Report` struct in `internal/analysis/analysis.go`.
    -   Add fields: `AvgBarsOverall`, `AvgBars1h`, `AvgSignalHealth`, `Has1hData`.

## Phase 2: Analysis Loop Update
- [x] Task: Modify `Analyze` function in `internal/analysis/analysis.go`.
    -   During the main loop:
        -   Accumulate sums for Overall Bars.
        -   Calculate Health Score for each sample and accumulate sum.
    -   After the loop (using `data` slice):
        -   Define `oneHourAgo = report.EndTime.Add(-1 * time.Hour)`.
        -   Iterate backwards from end of `data`.
        -   Accumulate bars for samples where `Time >= oneHourAgo`.
        -   Calculate final averages.

## Phase 3: Reporting & Verification
- [x] Task: Update `printReport` in `internal/analysis/analysis.go` to render the new "BARS AVG" section.
    -   Hide "Last 1h" if `!Has1hData`.
- [x] Task: Update `TestAnalyzeEndToEnd` in `analysis_test.go` to verify the new section exists and calculations are correct.
- [x] Task: Verify: Run `analyze` on `stats.log` and check the output against expected values.

## Phase 4: Chart Integration
- [x] Task: Export `CalculateSignalHealth` from `internal/analysis` (ensure it is public).
- [x] Task: Update `Generate` in `internal/charting/charting.go`.
    -   Calculate Signal Health for each data point using the exported function.
    -   Add `healthXYs` series to the "Signal Bars" plot.
    -   Style the new line (e.g., Green/Teal) and add a legend entry.
- [x] Task: Verify: Generate a chart and confirm both lines appear in the 3rd subplot.
