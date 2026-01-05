# Track Plan: Always-on JSON Logging

## Phase 1: Setup
- [x] Task: specific `stats.log` to `.gitignore`. [git-hash]
- [x] Task: Verify: Check that `stats.log` is ignored by git. [git-hash]

## Phase 2: Implementation
- [x] Task: Modify `main.go` to initialize a secondary `JSONLogger` targeting `stats.log`. [git-hash]
- [x] Task: Ensure this logger is set to append mode (which `NewJSONLogger` should already support, but verify). [git-hash]
- [x] Task: Update the main loop to write to this background logger on every tick. [git-hash]
- [x] Task: Handle potential file locking/collision if `-output stats.log` is passed (simple check: if user output == "stats.log", don't double write). [git-hash]

## Phase 3: Verification
- [x] Task: Run the app without flags. Verify `stats.log` is created and populated. [git-hash]
- [x] Task: Run the app with `-format csv -output test.csv`. Verify `stats.log` AND `test.csv` are populated. [git-hash]

