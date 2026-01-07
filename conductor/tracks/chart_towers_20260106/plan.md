# Track Plan: Add Tower ID to Connection Chart

## Phase 1: Logic & Mapping
- [x] Task: Update `Generate` in `internal/charting/charting.go`.
    -   Pass 1: Iterate data to collect unique `GNBID`s. Store in a slice/map. Sort them.
    -   Create a mapping function `getTowerY(gnbid int) float64`.
- [x] Task: Construct `towerXYs`.
    -   Iterate data, map GNBID to Y.
    -   Handle `0` (no signal).

## Phase 2: Plotting
- [x] Task: Rename Chart 4 to "Connection Info / Tower & Radio Bands".
- [x] Task: Add `lineTower` (Magenta, StepStyle).
- [x] Task: Construct custom `plot.ConstantTicks`.
    -   Add fixed Band ticks.
    -   Add dynamic Tower ticks starting at Y=5.
- [x] Task: Apply 24h smoothing logic to `towerXYs` if needed (same as `bandXYs`).

## Phase 3: Verification
- [x] Task: Verify: Generate chart. Confirm Bands are at bottom (Orange) and Towers are stacked above (Magenta).
