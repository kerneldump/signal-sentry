# Track Plan: Chart Grid Layout (2x2)

## Phase 1: Layout Calculation
- [x] Task: Update `Generate` function in `internal/charting/charting.go`.
    -   Change canvas size constants to `20 * vg.Inch` width and `8 * vg.Inch` height (or `16x10` if preferred).
    -   Define `colWidth` and `rowHeight`.
- [x] Task: Implement the 4 `draw.Canvas` rectangles using the new grid coordinates.
    -   **Lat/Loss:** `Min: (0, RowHeight)`, `Max: (ColWidth, Height)`
    -   **Signal:** `Min: (ColWidth, RowHeight)`, `Max: (Width, Height)`
    -   **Bars:** `Min: (0, 0)`, `Max: (ColWidth, RowHeight)`
    -   **Bands:** `Min: (ColWidth, 0)`, `Max: (Width, RowHeight)`

## Phase 2: Verification
- [x] Task: Verify: Generate a chart (`./tmobile-stats chart -output grid_test.png`) and confirm proper alignment and readability.
