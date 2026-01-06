# Track Spec: 24h Smoothing Trigger for All Charts

## Goal
Improve chart readability and generation performance for very long duration datasets (e.g., multi-day logs) by applying smoothing to all metrics, not just Signal Bars.

## Requirements

### 1. Trigger Conditions
*   **Signal Bars Chart:**
    *   Trigger: Duration > **2 hours**.
    *   Target Points: **300**.
*   **All Other Charts (Latency, RSRP/SINR, Bands):**
    *   Trigger: Duration > **24 hours**.
    *   Target Points: **600** (Higher fidelity for detailed metrics).

### 2. Implementation Logic
*   Modify `internal/charting/charting.go`.
*   Calculate `duration` once.
*   Define `shouldSmoothBars` (existing logic).
*   Define `shouldSmoothAll` (`duration > 24h` AND `len > 600`).
*   Apply `downsample` conditionally based on these flags.

### 3. Visual Handling
*   **Bands:** Downsampling discrete band data (1, 2, 3) results in floats (e.g., 2.5).
    *   *Decision:* When smoothed, the "Step Line" look might be weird. However, for 24h+ trends, seeing "Mostly Band X" is the goal. We will accept the averaging behavior for now as it reflects "mixed connection" during that bucket.

## Notes
*   This ensures "Daily" reports remain detailed, but "Weekly" reports become manageable trend lines.
