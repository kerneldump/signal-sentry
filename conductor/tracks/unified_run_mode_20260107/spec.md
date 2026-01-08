# Specification: Unified Run Mode (TUI + Web)

## Context
Currently, the application requires two separate processes to run the Live Monitor (TUI) and the Web Server (Charts). Users want the convenience of running a single command to start both.

## Requirements
1.  **New Flags:**
    *   `-web`: Boolean flag to enable the background web server.
    *   `-web-port`: Integer flag to specify the port (default: 8080).
2.  **Concurrency:**
    *   The web server must run in a separate goroutine so it doesn't block the TUI or the polling loop.
3.  **Logging Silence:**
    *   The web server's standard logging (e.g., "Starting server...", "Request received") must be suppressed or redirected when running in TUI mode to prevent screen corruption.
4.  **Integration:**
    *   The command `sudo go run . -live -web` should start the TUI *and* serve the web charts at `http://localhost:8080`.
    *   The web server will continue to read from `stats.log` (which the main process is writing to).

## User Experience
*   User runs: `signal-sentry -live -web`
*   The TUI appears as normal.
*   The user can open a browser to `localhost:8080` to see the charts.
*   No extra terminal output appears.
