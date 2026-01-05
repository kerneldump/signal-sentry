# Track Spec: Historical Analysis Tool

## Goal
Implement a subcommand `signal-sentry analyze` that parses the `stats.log` file and provides a statistical summary of the historical signal data.

## Requirements
1.  **CLI Command:**
    -   Detect `analyze` argument (or flag) to switch modes from "monitor" to "analyze".
    -   Default input file: `stats.log`.
    -   Optional flag: `-input <file>` to analyze a specific log file.
2.  **Analysis Metrics (Console Output):**
    -   **Time Range:** Start time - End time (Duration).
    -   **Total Samples:** Count of data points.
    -   **Signal Strength (RSRP):** Min, Max, Average.
    -   **Signal Quality (SINR):** Min, Max, Average.
    -   **Bands Seen:** List of unique bands (e.g., n41, n71) and % of time on each.
    -   **Towers Seen:** List of unique gNBIDs.
    -   **Ping Stats (if available):** Min/Avg/Max RTT, Packet Loss %.

## Implementation Details
-   Read `stats.log` line by line (streaming JSON decoder).
-   Accumulate stats in a `AnalysisReport` struct.
-   Print a formatted text table to stdout.

## Constraints
-   Must handle potentially large log files efficiently (streaming read).
-   Must handle malformed lines gracefully (skip and log warning).
