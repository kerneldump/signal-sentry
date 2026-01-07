# Track Spec: Improve Logarithmic Axis Labels

## Goal
Improve the readability of the "Latency & Packet Loss" chart by replacing scientific notation labels (e.g., `1e+01`, `1e+02`) with standard integer labels (e.g., `10`, `100`).

## Requirements

### 1. Custom Label Formatting
*   **Target:** Y-Axis of the "Latency & Packet Loss" chart.
*   **Format:** Standard decimal/integer representation.
*   **Scale:** Logarithmic.

### 2. Implementation
*   Override the default tick formatting for the logarithmic axis.
*   Ensure that major powers of 10 are labeled as `1`, `10`, `100`, `1000` etc.
*   If possible, provide enough intermediate ticks/grid lines to make values between powers of 10 (like 30ms or 50ms) easy to estimate.

## Verification
*   Generate a chart and confirm that Y-axis labels are integers.
