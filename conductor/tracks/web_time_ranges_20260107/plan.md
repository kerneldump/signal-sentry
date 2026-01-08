# Track Plan: Add Granular Time Ranges to Web UI

## Phase 1: Implementation
- [x] Task: Update `handleIndex` in `internal/web/server.go`.
    -   Insert `{"5m", "5m"}`, `{"15m", "15m"}`, `{"30m", "30m"}`, `{"45m", "45m"}` at the beginning of the `ranges` slice.
- [x] Task: Verify: Run `signal-sentry web` and request `/?range=5m` to ensure the page loads and the chart is generated without error. [Verified via unit test]
