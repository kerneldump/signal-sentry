# Track Spec: Ping Statistics & Advanced Configuration

## Goal
Enhance the monitoring capabilities by adding real-time network latency (ping) statistics and implementing a robust configuration management system to support new settings.

## Key Features

### 1. Configuration Management
-   **Config File:** Support loading settings from a file (JSON/YAML).
-   **Configurable Fields:**
    -   `RouterURL` (Signal API endpoint).
    -   `PingTarget` (IP/Domain to ping, e.g., "8.8.8.8").
    -   `RefreshInterval` (default for signal stats).
-   **Precedence:** CLI Flags > Config File > Defaults.
-   **New Flag:** `-config <path>` to specify config file.

### 2. Ping Statistics (Background Service)
-   **Metrics:**
    -   Minimum RTT
    -   Average RTT
    -   Maximum RTT
    -   Standard Deviation (StdDev)
    -   Packet Loss %
-   **Method:** Execute system `ping` command (shell out).
    -   Must support macOS (`ping` output format).
-   **Execution:**
    -   Asynchronous background process.
    -   Frequency: Fixed at 1 second (independent of signal fetch).
    -   Timeout: 1 second per ping.

### 3. Data & Output Updates
-   **Composite Data Model:** Combine `GatewayResponse` and `PingStats` into a unified `SentryState` or similar structure for logging/display.
-   **Output Channels:**
    -   **CLI/UI:** Add 5 new columns for the metrics.
    -   **CSV:** Append new columns.
    -   **JSON:** Add new fields.

## Constraints
-   Ping must not block the main application loop.
-   Must parse `ping` output reliably on Darwin (macOS).
