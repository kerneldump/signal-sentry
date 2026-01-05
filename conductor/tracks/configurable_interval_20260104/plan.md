# Track Plan: Add configurable refresh interval via command-line flag

## Phase 1: CLI Flag Implementation [checkpoint: 0802aba]
- [x] Task: Write failing tests for `-interval` flag validation (e.g., negative/zero values) [4af5792]
- [x] Task: Implement flag parsing and validation in `main.go` [2a58879]
- [x] Task: Conductor - User Manual Verification 'Phase 1: CLI Flag Implementation' (Protocol in workflow.md) [0802aba]

## Phase 2: Integration with Polling Loop [checkpoint: 88196f2]
- [x] Task: Refactor `main.go` to use the variable interval instead of the `refreshRate` constant [8e9d4d6]
- [x] Task: Verify the tool refreshes at the specified rate (manual and automated check) [8e9d4d6]
- [x] Task: Conductor - User Manual Verification 'Phase 2: Integration with Polling Loop' (Protocol in workflow.md) [88196f2]
