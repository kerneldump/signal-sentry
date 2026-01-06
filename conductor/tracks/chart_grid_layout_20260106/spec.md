# Track Spec: Chart Grid Layout (2x2)

## Goal
Optimize the visual layout of the generated charts by moving from a tall vertical stack (1x4) to a balanced grid (2x2). This improves readability on standard screens and reduces vertical scrolling.

## Requirements

### Layout Configuration
*   **Grid Dimensions:** 2 Columns x 2 Rows.
*   **Total Charts:** 4.
*   **Canvas Size:** Needs adjustment to accommodate side-by-side plots (e.g., wider and shorter than before).
    *   *Previous:* 10 inch Width x 16 inch Height (approx 4 inch per row).
    *   *Proposed:* 20 inch Width x 8 inch Height (maintains 10x4 aspect per plot).

### Plot Placement
1.  **Top-Left:** Latency & Packet Loss
2.  **Top-Right:** Signal Strength (RSRP/SINR)
3.  **Bottom-Left:** Signal Bars (Reported/Health)
4.  **Bottom-Right:** 5G Band

### Styling
*   Ensure fonts and axis labels remain legible with the new dimensions.
*   Maintain the existing color schemes and plot logic.

## Implementation Details
*   Modify `internal/charting/charting.go`.
*   Update `vgimg.NewWith` dimensions.
*   Recalculate `draw.Canvas` rectangles:
    *   `ColWidth = TotalWidth / 2`
    *   `RowHeight = TotalHeight / 2`
    *   Define `Min` and `Max` points for 4 distinct rectangles based on (X,Y) offsets.
