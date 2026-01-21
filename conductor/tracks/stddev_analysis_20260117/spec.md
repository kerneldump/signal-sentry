# Specification: StdDev/Jitter Analysis

## Goal
Add a new row "StdDev (ms)" to the "HISTORICAL SIGNAL ANALYSIS" output table, showing the Min, Avg, and Max of the ping standard deviation values collected in the samples.

## Requirements
1.  Modify `internal/analysis/analysis.go`:
    -   Update `Report` struct to include a `StdDev` metric (min, avg, max, sum, count).
    -   Update `Analyze` function to populate this metric from `stats.Ping.StdDev`.
    -   Update `printReport` function to output the new row under "Ping (ms)".
2.  Output Format:
    -   Label: "StdDev (ms)"
    -   Columns: Min, Avg, Max
    -   Formatting: Match the style of existing metrics (likely `%.1f`).

## Data Source
-   `models.CombinedStats.Ping.StdDev` (already exists in `internal/models/models.go`).
