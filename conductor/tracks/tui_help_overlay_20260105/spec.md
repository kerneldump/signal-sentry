# Track Spec: TUI Help Overlay

## Goal
Implement a togglable "Help/Info" overlay in the interactive TUI mode (`-live`) that displays the detailed signal metrics guide without stopping the background data collection.

## Requirements
1.  **Toggle Key:** Pressing `i` should toggle the view between the main signal monitor and the metrics guide.
2.  **Visual Layout:**
    *   **Header:** The Device Info, RSRP/SINR summary, and Ping Stats should remain visible at the top (to indicate the app is still running).
    *   **Content Area:** The scrolling signal history is replaced by the static "SIGNAL METRICS GUIDE" text when active.
    *   **Footer/Instruction:** Update the interval line to include "i for info" or similar.
3.  **Data Continuity:** Background fetching (Gateway + Ping) must continue uninterrupted. New data points are added to the history buffer even while the user is viewing the help screen.

## Content
The text to display matches the `printLegend()` output from `main.go`.

```text
SIGNAL METRICS GUIDE:
---------------------
* BAND:  The frequency band in use.
         n41: High speed, shorter range (Ultra Capacity).
         n25: Balanced speed and range.
         n71: Long range, slower speeds.

* RSRP:  (Reference Signal Received Power) Your main signal strength.
         Excellent > -80  | Good -80 to -95
         Fair -95 to -110 | Poor < -110 (Risk of drops).

* SINR:  (Signal-to-Interference-plus-Noise Ratio) Signal quality.
         Higher is better. > 20 is excellent.
         < 0 means high noise (your speed will suffer).

* RSRQ:  (Reference Signal Received Quality) The congestion indicator.
         If SINR is Good (high) but RSRQ is Bad (low), the tower is
         likely congested with heavy traffic.

* CID & gNBID:
         gNBID identifies the physical TOWER.
         CID identifies the SECTOR (which side of the tower you are facing).
```
