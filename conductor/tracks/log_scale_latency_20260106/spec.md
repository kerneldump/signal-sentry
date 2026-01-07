# Track Spec: Logarithmic Scale for Latency Chart

## Goal
Improve the readability of the "Latency & Packet Loss" chart. Currently, high outliers (spikes in StdDev or Latency) cause the Y-axis to auto-scale too high, making the typical latency (e.g., 20-50ms) appear as a flat line.

## Requirements

### 1. Chart Configuration
*   **Target:** Chart 1 (Latency & Packet Loss).
*   **Scale:** Change Y-Axis from Linear to **Logarithmic**.
*   **Range:**
    *   **Min:** Must be > 0 (Log(0) is undefined). Set nice minimum (e.g., 1ms or 10ms).
    *   **Max:** Auto-scale to fit data.

### 2. Data Handling
*   **Zero Values:** Packet loss is often `0.0`. Log(0) crashes or errors.
    *   **Solution:** For plotting, clamp min values to a small epsilon (e.g., 0.1) or strictly use Linear scale for Loss if plotted on separate axis.
    *   *Decision:* Since Loss is on the same chart as Latency (ms), and Loss is 0-100%, while Latency is 20-500ms...
        *   Log scale works for Latency (20 to 500).
        *   Log scale for Loss (0 to 100) is weird if Loss is 0.
    *   **Alternative:** The user complained about *StdDev* spikes.
    *   We will apply Log Scale to the Y-axis. We must ensure `0` values (Loss) are handled (e.g., mapped to 0.1 or not plotted if we want to hide them, or just accept that Log scale visualizes "0" effectively at the bottom if min is set low).

## Implementation Details
*   Modify `internal/charting/charting.go`.
*   Set `pLat.Y.Scale = plot.LogScale{}`.
*   Update Y-Ticks to be friendlier for Log scale (e.g., 10, 100, 1000).
*   **Safety:** Iterate data points. If Y <= 0, set to 0.1 (or similar non-zero positive) to prevent panic.

## Verification
*   Generate chart with a known spike. Verify the "normal" baseline is still readable.
