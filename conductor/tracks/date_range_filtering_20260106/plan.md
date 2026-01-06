# Track Plan: Date Range Filtering for Analysis

## Phase 1: Core Logic & Flag Parsing
- [x] Task: Create `internal/analysis/filter.go` (or similar).
    -   Define `TimeFilter` struct holding Start and End times.
    -   Implement `NewTimeFilter(start, end string, rangeDur time.Duration) (*TimeFilter, error)` to handle the parsing logic and precedence rules.
    -   Implement `ParseISOTime(s string) (time.Time, error)` helper for YYYY-MM-DD [HH:MM:SS].
- [x] Task: Write unit tests for `NewTimeFilter` to verify precedence (e.g., `-range` vs `-start`) and parsing.

## Phase 2: Analysis Integration
- [x] Task: Update `Analyze` and `ParseLog` signatures in `internal/analysis/analysis.go`.
    -   Option A: Pass `TimeFilter` to `Analyze`.
    -   Option B: Pass `TimeFilter` to `ParseLog` (more efficient, avoids loading unwanted data into memory). *Decision: Pass to ParseLog.*
- [x] Task: Modify `ParseLog` loop to check `filter.Contains(entry.Time)`.
- [x] Task: Verify: Update existing tests in `analysis_test.go` to accommodate signature changes (pass nil/empty filter).

## Phase 3: CLI Wiring
- [x] Task: Update `runAnalysis` in `main.go`.
    -   Add flags: `-start`, `-end`, `-range`.
    -   Parse flags into `TimeFilter`.
    -   Pass filter to `analysis.Run`.
- [x] Task: Update `runChart` in `main.go`.
    -   Add flags: `-start`, `-end`, `-range`.
    -   Parse flags and pass to `analysis.ParseLog`.

## Phase 4: Output Polish
- [x] Task: Update `PrintReport` in `internal/analysis` to display the active filter (if any) in the header.
- [x] Task: Verify: Run `analyze -range 1h` and ensure only recent data is shown.
- [x] Task: Verify: Run `chart -range 24h` and ensure the graph covers only that period.
