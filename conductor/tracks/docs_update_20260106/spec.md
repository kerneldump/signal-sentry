# Track Spec: Documentation Update

## Goal
Update `README.md` to accurately reflect the tool's current capabilities, specifically focusing on the new Analysis and Charting features added in recent tracks.

## Requirements

### 1. Update `Usage` Section
*   **Filtering:** Document the new `-range`, `-start`, and `-end` flags under both `analyze` and `chart` subcommands. Provide examples.
    *   Example: `signal-sentry analyze -range 24h`
    *   Example: `signal-sentry chart -start 2026-01-01 -end 2026-01-02`

### 2. Update `Features` Section
*   **Signal Health:** Explain the new "Signal Health" metric (RSRP+SINR derived score).
*   **Smart Charting:** Mention the 2x2 grid layout and automatic smoothing for long durations.

### 3. Update `Analysis Report` Example
*   Replace the old sample output with a new one showing:
    *   Filter headers.
    *   "live" / "real-time" markers.
    *   `BARS AVG` section.

## Implementation Details
*   Modify `README.md` in place.
