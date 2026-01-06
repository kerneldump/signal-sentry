# Track Plan: Chart Smoothing for Signal Bars

## Phase 1: Logic Implementation
- [x] Task: Create `downsample` function in `internal/charting/charting.go`.
    -   Input: `plotter.XYs`, `maxPoints int`.
    -   Output: `plotter.XYs` (averaged).
    -   Logic: Iterate through input in chunks of `stride`. Average Y. Use Average X (or Start X).

## Phase 2: Integration & Styling
- [x] Task: Update `Generate` in `internal/charting/charting.go`.
    -   Add logic to detect if downsampling is needed (`len > 300`).
    -   Apply downsampling to `barsXYs` and `healthXYs` independently.
    -   Conditionally set `StepStyle`:
        -   If raw: `plotter.PreStep`.
        -   If smoothed: `plotter.NoStep` (default line).

## Phase 3: Verification
- [x] Task: Verify: Generate a chart for a long duration (full log) and ensure it looks smooth.
- [x] Task: Verify: Generate a chart for a short duration (`-range 30m`) and ensure it remains raw/stepped.
