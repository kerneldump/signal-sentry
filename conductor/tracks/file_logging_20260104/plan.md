# Track Plan: Add JSON and CSV file logging support

## Phase 1: Flag Parsing and Logger Scaffolding [checkpoint: 4210f78]
- [x] Task: Update `main.go` to support `-format` and `-output` flags (requires coordination with the interval track) [cb9365b]
- [x] Task: Create `internal/logger` package with `Logger` interface [eec85b8]
- [x] Task: Conductor - User Manual Verification 'Phase 1: Flag Parsing and Logger Scaffolding' (Protocol in workflow.md) [4210f78]

## Phase 2: JSON Implementation [checkpoint: fd470e7]
- [x] Task: Implement `JSONLogger` in `internal/logger` [9ea7de8]
- [x] Task: Write tests for JSON logging (verifying file output) [9ea7de8]
- [x] Task: Integration: Hook up JSON logging in `main.go` [27a019e]
- [x] Task: Conductor - User Manual Verification 'Phase 2: JSON Implementation' (Protocol in workflow.md) [fd470e7]

## Phase 3: CSV Implementation
- [x] Task: Implement `CSVLogger` in `internal/logger` (handling headers and row mapping) [32557eb]
- [x] Task: Write tests for CSV logging [32557eb]
- [x] Task: Integration: Hook up CSV logging in `main.go` [c94a4d6]
- [ ] Task: Conductor - User Manual Verification 'Phase 3: CSV Implementation' (Protocol in workflow.md)
