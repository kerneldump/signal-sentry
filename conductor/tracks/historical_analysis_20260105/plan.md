# Track Plan: Historical Analysis Tool

## Phase 1: Subcommand Structure
- [ ] Task: Refactor `main.go` to handle subcommands (currently it just parses flags for the main loop).
    -   Check `os.Args` for "analyze".
    -   If found, delegate to `internal/analysis.Run()`.
    -   Else, proceed with existing monitoring logic.

## Phase 2: Analysis Logic
- [ ] Task: Create `internal/analysis` package.
- [ ] Task: Implement `AnalyzeFile(filepath string) (*Report, error)` function.
- [ ] Task: Implement streaming JSON decoder to read `GatewayResponse` objects (and `PingStats` if merged).
- [ ] Task: Calculate Min/Max/Avg for RSRP, SINR, Bars.
- [ ] Task: Track frequency maps for Bands and Towers.

## Phase 3: Reporting
- [ ] Task: Implement `PrintReport(r *Report)` to format the stats nicely (using tabwriter or formatted printf).
- [ ] Task: Verify: Create a sample `stats.log` with known values. Run `signal-sentry analyze` and check the math.
