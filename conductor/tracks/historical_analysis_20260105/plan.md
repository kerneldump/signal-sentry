# Track Plan: Historical Analysis Tool

## Phase 1: Subcommand Structure
- [x] Task: Refactor `main.go` to handle subcommands (currently it just parses flags for the main loop). [git-hash]
    -   Check `os.Args` for "analyze".
    -   If found, delegate to `internal/analysis.Run()`.
    -   Else, proceed with existing monitoring logic.

## Phase 2: Analysis Logic
- [x] Task: Create `internal/analysis` package. [git-hash]
- [x] Task: Implement `AnalyzeFile(filepath string) (*Report, error)` function. [git-hash]
- [x] Task: Implement streaming JSON decoder to read `GatewayResponse` objects (and `PingStats` if merged). [git-hash]
- [x] Task: Calculate Min/Max/Avg for RSRP, SINR, Bars. [git-hash]
- [x] Task: Track frequency maps for Bands and Towers. [git-hash]

## Phase 3: Reporting
- [x] Task: Implement `PrintReport(r *Report)` to format the stats nicely (using tabwriter or formatted printf). [git-hash]
- [x] Task: Enhance report with live tower identification and global loss calculation. [873620c]
- [x] Task: Verify: Create a sample `stats.log` with known values. Run `signal-sentry analyze` and check the math. [git-hash]

