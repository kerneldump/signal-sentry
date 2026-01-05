# Track Plan: Add JSON and CSV file logging support

## Phase 1: Flag Parsing and Logger Scaffolding
- [ ] Task: Update `main.go` to support `-format` and `-output` flags (requires coordination with the interval track)
- [ ] Task: Create `internal/logger` package with `Logger` interface
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Flag Parsing and Logger Scaffolding' (Protocol in workflow.md)

## Phase 2: JSON Implementation
- [ ] Task: Implement `JSONLogger` in `internal/logger`
- [ ] Task: Write tests for JSON logging (verifying file output)
- [ ] Task: Integration: Hook up JSON logging in `main.go`
- [ ] Task: Conductor - User Manual Verification 'Phase 2: JSON Implementation' (Protocol in workflow.md)

## Phase 3: CSV Implementation
- [ ] Task: Implement `CSVLogger` in `internal/logger` (handling headers and row mapping)
- [ ] Task: Write tests for CSV logging
- [ ] Task: Integration: Hook up CSV logging in `main.go`
- [ ] Task: Conductor - User Manual Verification 'Phase 3: CSV Implementation' (Protocol in workflow.md)
