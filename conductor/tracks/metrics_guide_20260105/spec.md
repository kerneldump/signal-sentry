# Track Spec: Update Signal Metrics Guide

## Goal
Update the embedded signal metrics documentation within the application to provide more accurate thresholds and descriptions based on real-world 5G behavior (specifically for n25, n41, and n71 bands).

## Requirements
1.  **Adjust RSRP Thresholds:**
    -   Excellent: > -80
    -   Good: -80 to -95
    -   Fair: -95 to -110
    -   Poor: < -110
2.  **Update Band Descriptions:**
    -   Include **n25** (Balanced speed and range).
    -   Refine n41 (High speed/Ultra Capacity) and n71 (Long range).
3.  **Refine Identification Labels:**
    -   Clarify that `gNBID` is the physical tower and `CID` is the specific sector.
4.  **Consistency:**
    -   Ensure this guide is displayed correctly in the CLI/TUI modes.

## Content to Implement
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
