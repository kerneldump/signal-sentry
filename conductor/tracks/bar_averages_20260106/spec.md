# Track Spec: Bar Averages and Derived Signal Metrics

## Goal
Enhance the `analyze` report to provide more granular insight into signal "Bars". This includes calculating the decimal average of reported bars and introducing a "Signal Health" metric that combines RSRP and SINR to represent true connection quality.

## Requirements

### 1. New Report Section: `BARS AVG`
A new section will be added to the bottom of the `analyze` output:

```text
BARS AVG:
Overall      3.5
Last 1h      3.2
SgnlHealth   3.4
```

### 2. Metric Definitions

*   **Overall:**
    *   The arithmetic mean of the `bars` reported by the gateway for **all samples** included in the current report (respecting any `-range` or `-start`/`-end` filters applied).
    *   Format: 1 decimal place (e.g., `3.5`).

*   **Last 1h:**
    *   The arithmetic mean of the `bars` reported by the gateway for samples within the **last 1 hour** of the *available data range*.
    *   **Logic:**
        *   Determine `DataEndTime`.
        *   Filter samples where `Time >= DataEndTime - 1h`.
        *   Calculate average.
    *   **Visibility:** This row is **hidden** if the total duration of the data in the report is less than 1 hour.

*   **SgnlHealth (Derived Quality Score):**
    *   A calculated score (1.0 - 5.0) representing the "True Quality" of the connection for each sample, averaged over the report period.
    *   **Formula:** `Score = (RSRP_Score * 0.7) + (SINR_Score * 0.3)`
    *   **RSRP Scoring (Strength):**
        *   `> -80`   = 5.0
        *   `-80 to -90`  = 4.0
        *   `-90 to -100` = 3.0
        *   `-100 to -110`= 2.0
        *   `< -110`  = 1.0
    *   **SINR Scoring (Quality):**
        *   `> 20`    = 5.0
        *   `13 to 20`    = 4.0
        *   `0 to 13`     = 3.0
        *   `-10 to 0`    = 2.0
        *   `< -10`   = 1.0
    *   *Note: SINR thresholds derived from standard LTE/5G performance tiers.*

## Implementation Details

### Data Structures
*   Update `analysis.Report` struct to hold:
    *   `AvgBarsOverall float64`
    *   `AvgBars1h float64`
    *   `AvgSignalHealth float64`
    *   `Has1hData bool`

### Analysis Logic
*   During the iteration of log entries:
    1.  Accumulate `SumBars` and `Count` for "Overall".
    2.  Calculate `DerivedScore` for the current sample using the formula. Accumulate `SumHealth` and `Count`.
    3.  Store samples (or just necessary values + timestamp) in a buffer/slice to calculate "Last 1h" *after* the loop (since we need to know the final `EndTime` to define the 1h window), OR do a second pass/smart accumulation.
        *   *Decision:* Since we already parse into a slice in memory (`[]models.CombinedStats`), we can easily iterate backwards from the end to calculate the "Last 1h" average efficiently.

### Output
*   Update `printReport` to render the `BARS AVG` table using `tabwriter` for alignment.
