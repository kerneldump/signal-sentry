# Track Spec: Add Tower ID to Connection Chart

## Goal
Visualize Tower ID changes alongside 5G Band changes to help users correlate tower hopping with band switching. Instead of adding a 5th chart, we will combine these into a single "Connection Info" chart.

## Requirements

### 1. Visualization
*   **Chart:** Reuse Chart 4 (Bottom-Right). Rename from "5G Band" to "Connection Info".
*   **Series 1 (Bands):** Orange Step Line (Y-levels 1-3).
    *   1: n71
    *   2: n25
    *   3: n41
*   **Series 2 (Towers):** Magenta Step Line (Y-levels 5+).
    *   Dynamically map unique Tower IDs found in the data to integers 5, 6, 7...
    *   Map `0` (No Signal) to `0`.

### 2. Y-Axis
*   **Custom Ticks:** The Y-axis must display text labels for both bands and towers.
    *   `1` -> "n71"
    *   `2` -> "n25"
    *   `3` -> "n41"
    *   `5` -> "Tower 123..."
    *   `6` -> "Tower 456..."

### 3. Logic
*   Scan all data points first to collect unique `GNBID` values.
*   Sort Unique Towers (e.g., ascending ID).
*   Assign offsets.
*   Generate `towerXYs` points.
*   Apply **Smoothing** if applicable (using the same `shouldSmoothAll` or `shouldSmoothBars` trigger? Probably `shouldSmoothAll` since it's discrete data like bands). *Decision: Use raw step lines unless the "All" trigger (>24h) is hit, similar to Bands.*

## Implementation Details
*   Modify `internal/charting/charting.go`.
*   Maintain the 2x2 Grid layout.
