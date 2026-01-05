# Track Plan: Interactive Live View (CLI)

## Phase 1: Dependencies and Core UI Structure
- [ ] Task: Initialize `internal/ui` package and add `bubbletea` dependency (`go get github.com/charmbracelet/bubbletea`).
- [ ] Task: Define the `Model` struct in `internal/ui` (holding the buffer, current interval, window size).
- [ ] Task: Implement the `Init` and basic `View` methods (placeholder display).
- [ ] Task: Implement `Update` to handle `tea.WindowSizeMsg` and `tea.KeyMsg` (quit logic).
- [ ] Task: Verify: Create a simple main entry point (temporary) to run the UI and check rendering/resizing.

## Phase 2: Data Integration and Logic
- [ ] Task: Define a `Msg` type for new signal data.
- [ ] Task: Create a command/tick loop in `bubbletea` to fetch data from `gateway` client.
- [ ] Task: Implement the "Rolling Buffer" logic (prepend new item, trim > 30).
- [ ] Task: Implement dynamic interval adjustment (handle +/- keys and update the tick duration).
- [ ] Task: Verify: Run the UI and ensure data populates and updates at the set interval.

## Phase 3: CLI Integration and Polish
- [ ] Task: Update `main.go` to parse `-live` flag.
- [ ] Task: Branch execution in `main.go`: if `-live` is true, start the Bubble Tea program; otherwise run the legacy loop.
- [ ] Task: Polish the `View` rendering (align columns, header, hide lines if screen is short).
- [ ] Task: Verify: Test `-live` mode and legacy mode to ensure no regression.
