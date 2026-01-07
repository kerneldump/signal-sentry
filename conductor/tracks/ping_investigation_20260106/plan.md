# Track Plan: Ping Reliability Investigation

## Phase 1: Code Audit & Analysis
- [x] Task: Delegate to `codebase_investigator` to analyze `internal/pinger/pinger.go` and its usage of `pro-bing`.
    -   Check initialization overhead.
    -   Check timeout configuration.
- [x] Task: Create a standalone test script `cmd/pingtest/main.go` to compare "One-Shot" (current) vs "Continuous" ping behavior.

## Phase 2: Refactoring (If confirmed)
- [x] Task: Refactor `Pinger` to use a persistent `pro-bing` instance running in the background.
    -   *Note: Implemented "One-Shot" with decoupled timeout instead of full continuous mode to resolve the primary issue (false positives) with less complexity.*
    -   Use `pinger.Run()` (blocking) in a goroutine.
    -   Use callbacks `OnRecv`, `OnFinish` (not needed for continuous), or just `OnRecv` to update stats.
    -   Calculate stats (Min/Max/Avg) manually using Welford's algorithm or a sliding window since `pro-bing` stats are cumulative for the whole run.
- [x] Task: Decouple `Timeout` from `Interval`. Set Timeout to something safe (e.g. 1000ms or 2000ms) while keeping the send interval at 1s.

## Phase 3: Verification
- [x] Task: Run the refactored pinger alongside the old one (or just standard ping) to verify stability.
