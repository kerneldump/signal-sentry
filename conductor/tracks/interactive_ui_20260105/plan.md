# Track Plan: Interactive Live View (CLI)

## Phase 1: Dependencies and Core UI Structure
- [x] Task: Initialize `internal/ui` package and add `bubbletea` dependency (`go get github.com/charmbracelet/bubbletea`). [git-hash]
- [x] Task: Define the `Model` struct in `internal/ui` (holding the buffer, current interval, window size). [git-hash]
- [x] Task: Implement the `Init` and basic `View` methods (placeholder display). [git-hash]
- [x] Task: Implement `Update` to handle `tea.WindowSizeMsg` and `tea.KeyMsg` (quit logic). [git-hash]
- [x] Task: Verify: Create a simple main entry point (temporary) to run the UI and check rendering/resizing. [git-hash]

## Phase 2: Data Integration and Logic
- [x] Task: Define a `Msg` type for new signal data. [git-hash]
- [x] Task: Create a command/tick loop in `bubbletea` to fetch data from `gateway` client. [git-hash]
- [x] Task: Implement the "Rolling Buffer" logic (prepend new item, trim > 30). [git-hash]
- [x] Task: Implement dynamic interval adjustment (handle +/- keys and update the tick duration). [git-hash]
- [x] Task: Verify: Run the UI and ensure data populates and updates at the set interval. [git-hash]

## Phase 3: CLI Integration and Polish
- [x] Task: Update `main.go` to parse `-live` flag. [git-hash]
- [x] Task: Branch execution in `main.go`: if `-live` is true, start the Bubble Tea program; otherwise run the legacy loop. [git-hash]
- [x] Task: Polish the `View` rendering (align columns, header, hide lines if screen is short). [git-hash]
- [x] Task: Verify: Test `-live mode and legacy mode to ensure no regression. [git-hash]

