# Specification: Display Averages on Latency Chart

## Context
Currently, the "Latency & Packet Loss" chart displays a label attached to the last data point showing the value of that specific point (e.g., "45ms"). This can be misleading as it represents a momentary snapshot rather than the overall performance during the selected time window.

## Goal
Update the chart to display the **average** value over the entire selected time range for both:
1.  **Ping Latency**
2.  **Standard Deviation (Jitter)**

## Requirements
1.  **Calculation:**
    -   Calculate the arithmetic mean of all valid (non-zero) Ping Latency values in the dataset.
    -   Calculate the arithmetic mean of all Standard Deviation values in the dataset.
2.  **Visualization:**
    -   Replace the existing "last point value" label.
    -   Attach the new label to the last data point (same position as before).
    -   Format the label as "Avg: [Value]".
    -   Display values with **1 decimal place** (e.g., "Avg: 27.3ms").
3.  **Scope:**
    -   Apply this change only to the "Latency & Packet Loss" chart.
