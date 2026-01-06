# Track Spec: Date Range Filtering for Analysis

## Goal
Enable users to filter the input data for `analyze` and `chart` subcommands based on specific time ranges. This allows for focused troubleshooting of recent issues or historical comparisons.

## Requirements

### 1. New Command-Line Flags
The following flags should be added to both `analyze` and `chart` subcommands:

*   `-start string`: The start time of the window.
    *   **Format:** ISO 8601 Date (`YYYY-MM-DD`) or Date+Time (`YYYY-MM-DD HH:MM:SS`).
    *   **Default:** Beginning of the log (if unset).
*   `-end string`: The end time of the window.
    *   **Format:** ISO 8601 Date (`YYYY-MM-DD`) or Date+Time (`YYYY-MM-DD HH:MM:SS`).
    *   **Default:** Now (if unset).
*   `-range duration`: A relative duration looking back from **NOW**.
    *   **Format:** Go duration string (e.g., `24h`, `30m`, `1h30m`).
    *   **Behavior:** If specified, it overrides `-start` and sets it to `Now() - duration`.

### 2. Logic & Precedence
1.  **Relative Range (`-range`):**
    *   If provided, `StartTime = time.Now().Add(-duration)` and `EndTime = time.Now()`.
    *   This takes precedence over `-start`.
2.  **Explicit Date Range (`-start` / `-end`):**
    *   Users can provide just `-start` (from date X until now).
    *   Users can provide just `-end` (from beginning of log until date Y).
    *   Users can provide both (from date X to date Y).
3.  **Filtering:**
    *   The application must filter log entries based on their `localTime` (or `upTime` derived timestamp).
    *   Only entries falling within `[StartTime, EndTime]` (inclusive) are processed for reports or charts.

### 3. User Interface
*   The `analyze` report header should reflect the *requested* time range vs the *actual* data found.
    *   *Current:* "Time Range: <FirstSample> to <LastSample>"
    *   *New:* "Filter: <Start> to <End>" (if filter active) + "Data Range: ..."

## Examples
```bash
# Analyze the last 24 hours
./tmobile-stats analyze -range 24h

# Chart a specific incident window
./tmobile-stats chart -start "2026-01-05 14:00:00" -end "2026-01-05 16:00:00"

# Compare today's stats
./tmobile-stats analyze -start "2026-01-06"
```
