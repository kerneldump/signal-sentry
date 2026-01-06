# Track Spec: Chart Smoothing for Signal Bars

## Goal
Improve the readability of the "Signal Bars" chart when visualizing large datasets (e.g., 24h+). Currently, the high-frequency switching between integer bar values (3/4) creates visual noise. We will implement dynamic aggregation to show smoothed trends instead.

## Requirements

### 1. Trigger Condition
*   Smoothing is applied if the total number of data points exceeds a target threshold (e.g., **300 points**).
*   If points <= 300, the chart remains in "Raw Mode" (Step Line, Full Resolution).

### 2. Aggregation Logic (Downsampling)
*   **Target:** Reduce dataset to ~300 points.
*   **Algorithm:** Simple Bucket Averaging.
    *   `BucketSize = TotalPoints / 300`.
    *   For each bucket, calculate the arithmetic mean of Y values.
    *   X value for the bucket can be the midpoint or start.
*   **Scope:** Applied **ONLY** to the Signal Bars chart (both "Reported Bars" and "Signal Health" series). Other charts remain raw.

### 3. Visual Changes (Only when Smoothed)
*   **Reported Bars (Black Line):**
    *   Change from `StepStyle` (Discrete) to **Normal Line** (Continuous).
    *   This better represents the "Average Bar Value" (e.g., 3.4) drifting over time.
*   **Signal Health (Grey Area):**
    *   No style change needed, just using the smoothed data points for the polygon.

## Implementation Details
*   Modify `internal/charting/charting.go`.
*   Implement a `downsample(data plotter.XYs, maxPoints int) plotter.XYs` helper function.
*   In `Generate`:
    *   Populate `barsXYs` and `healthXYs` with full data first.
    *   Check `len(barsXYs) > 300`.
    *   If true:
        *   `barsXYs = downsample(barsXYs, 300)`
        *   `healthXYs = downsample(healthXYs, 300)`
        *   Disable `StepStyle` for the black line.
