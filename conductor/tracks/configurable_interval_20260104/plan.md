# Track Plan: Add configurable refresh interval via command-line flag

## Phase 1: CLI Flag Implementation
- [x] Task: Write failing tests for `-interval` flag validation (e.g., negative/zero values) [4af5792]
- [~] Task: Implement flag parsing and validation in `main.go`
- [ ] Task: Conductor - User Manual Verification 'Phase 1: CLI Flag Implementation' (Protocol in workflow.md)

## Phase 2: Integration with Polling Loop
- [ ] Task: Refactor `main.go` to use the variable interval instead of the `refreshRate` constant
- [ ] Task: Verify the tool refreshes at the specified rate (manual and automated check)
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Integration with Polling Loop' (Protocol in workflow.md)
