# Specification: Add Granular Time Ranges to Web UI

## Context
The current Web UI offers time ranges starting from 1 hour ("1h"). Users need more granular control to view immediate short-term trends, especially when diagnosing active issues or testing antenna adjustments.

## Requirements
1.  Add the following time range options to the navigation bar in the Web UI:
    *   **5m** (5 Minutes)
    *   **15m** (15 Minutes)
    *   **30m** (30 Minutes)
    *   **45m** (45 Minutes)
2.  These options should precede the existing "1h" option.
3.  The backend `handleChart` logic must support these new values (which it should already via `time.ParseDuration`).

## User Experience
*   The top navigation bar will now start with "5m | 15m | 30m | 45m | 1h ...".
*   Clicking these links will reload the page with `?range=Xm` and update the chart accordingly.
