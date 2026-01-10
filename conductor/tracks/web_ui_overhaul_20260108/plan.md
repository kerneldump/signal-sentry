# Track Plan: Web UI Time Range & Filtering Overhaul

## Phase 1: Backend Parsing Logic [checkpoint: 9656177]
- [x] Task: Create tests for `internal/web` to verify parsing of new query parameters (`range`, `start`, `end`). f6c8e04
- [x] Task: Update `internal/web/server.go` (or `request_handler.go`) to implement the parameter parsing logic. f6c8e04
- [x] Task: Conductor - User Manual Verification 'Backend Parsing Logic' (Protocol in workflow.md) f6c8e04

## Phase 2: Template Refactoring & Custom Input
- [ ] Task: Update the Go templates (HTML) for the toolbar.
    -   Remove old preset buttons.
    -   Add new preset buttons: `10m`, `1h`, `6h`, `24h`, `Max`.
    -   Add `<input type="text">` for custom relative duration.
    -   Add `<div>` container for date range pickers (initially hidden or alongside).
    -   Add `<input type="datetime-local">` for start and end times.
- [ ] Task: Conductor - User Manual Verification 'Template Refactoring & Custom Input' (Protocol in workflow.md)

## Phase 3: Frontend Interaction (JavaScript)
- [ ] Task: Implement JavaScript logic for "Auto-Update".
    -   Attach event listeners to preset buttons.
    -   Attach `change` (or `input` with debounce) listeners to custom text and date inputs.
    -   Implement logic to construct the new URL and reload the page (or fetch via AJAX if existing architecture supports it; assuming page reload for now based on simplicity).
- [ ] Task: Implement "Mutual Exclusivity" logic.
    -   If a preset or custom duration is touched, clear the Date Range inputs.
    -   If a Date Range input is touched, clear the custom duration input/active preset.
- [ ] Task: Conductor - User Manual Verification 'Frontend Interaction (JavaScript)' (Protocol in workflow.md)

## Phase 4: Integration & Polish
- [ ] Task: Verify the entire flow.
    -   Check if `analysis.Run` (or the chart generation function) respects the new high-precision filters.
    -   Ensure chart titles/subtitles reflect the active filter.
- [ ] Task: Polish the CSS (Bootstrap) to ensure the new inputs look good on mobile and desktop.
- [ ] Task: Conductor - User Manual Verification 'Integration & Polish' (Protocol in workflow.md)
