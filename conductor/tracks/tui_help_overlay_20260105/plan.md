# Track Plan: TUI Help Overlay

## Phase 1: Model Update
- [x] Task: Update `internal/ui/model.go` struct to include `showHelp bool`. [git-hash]
- [x] Task: Define the constant `helpText` string containing the metrics guide. [git-hash]

## Phase 2: Logic Implementation
- [x] Task: Update `Update` method in `internal/ui/model.go` to handle `KeyMsg("i")`: toggle `m.showHelp`. [git-hash]
- [x] Task: Update the instructions line in the View to say "(Press +/- to adjust, i for info, q to quit)". [git-hash]

## Phase 3: View Rendering
- [x] Task: Refactor `View` method. [git-hash]
    -   Always render Header (Device + Ping Stats).
    -   If `m.showHelp` is true: Render the static `helpText`.
    -   Else: Render the Table Header + Rolling Buffer (existing logic).
- [x] Task: Verify: Run `-live` mode, toggle `i`, check if background data (ping stats in header) keeps updating. [git-hash]

