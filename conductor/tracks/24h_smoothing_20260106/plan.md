# Track Plan: 24h Smoothing Trigger for All Charts

## Phase 1: Logic Update
- [x] Task: Update `Generate` in `internal/charting/charting.go`.
    -   Calculate `duration` in seconds.
    -   Define `shouldSmoothBars = duration > 2h`.
    -   Define `shouldSmoothAll = duration > 24h`.
- [x] Task: Apply downsampling:
    -   If `shouldSmoothBars`: Downsample Bars/Health to 300 pts.
    -   If `shouldSmoothAll`: Downsample Latency, RSRP, SINR, Bands to 600 pts.
    -   *Optimization:* If `shouldSmoothAll` is true, it implies `shouldSmoothBars` is also true (24h > 2h). We might want to use 600 pts for bars too in that case for consistency? *Decision: Keep bars at 300 for smoothness, others at 600 for detail.*

## Phase 2: Verification
- [x] Task: Verify: Generate chart for < 2h (Raw).
- [x] Task: Verify: Generate chart for 12h (Bars smooth, others raw).
- [x] Task: Verify: Generate chart for 30h (All smooth).
