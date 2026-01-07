# Specification: Lightweight Web Server for Live Charts

## Goal
Create a standalone subcommand (`web`) that launches a local HTTP server. This server will render the 2x2 signal stats chart (generated via `internal/charting`) and serve it within a simple HTML dashboard.

## Requirements

### 1. Subcommand
- **Command:** `signal-sentry web [flags]`
- **Flags:**
    - `-port <int>`: Port to listen on (default: `8080`).
    - `-input <string>`: Path to the log file (default: `stats.log`).
    - `-dev`: (Optional) Enable hot-reloading or verbose logging if needed.

### 2. Routes
- **`GET /`**: Serves the main HTML dashboard.
    - Accepts query parameter `range` (e.g., `?range=6h`). Default: `24h`.
    - Includes navigation links to switch ranges.
    - Embeds the chart image using `<img src="/chart.png?range=...">`.
    - Includes a mechanism to auto-refresh the page every 60 seconds.

- **`GET /chart.png`**: Generates and returns the PNG image.
    - Accepts query parameter `range` (string, duration format like `1h`, `24h` or `max`).
    - Parses the log file defined by `-input` using the specified range filter.
    - Generates the chart in-memory (buffer) or to a temp file using `internal/charting`.
    - Returns the image with correct Content-Type `image/png`.

### 3. User Interface (HTML)
- **Header:** Title ("Signal Sentry Live"), Current Range.
- **Navigation:** Links for: `1h`, `2h`, `3h`, `6h`, `12h`, `24h`, `48h`, `Max`.
- **Content:** The generated 2x2 chart image, centered.
- **Footer:** Last updated timestamp.
- **Behavior:**
    - Auto-refresh every 60 seconds.
    - Changing the range updates the URL and immediately reloads the chart.

### 4. Technical Constraints
- Reuse existing `internal/analysis` for parsing.
- Reuse existing `internal/charting` for image generation.
- **Crucial:** `internal/charting.Generate` currently writes to a file. It may need refactoring or wrapping to write to an `io.Writer` (like `http.ResponseWriter`) to avoid disk I/O churn, OR we can write to a temporary file and serve that. *Decision: Refactor `Generate` to accept `io.Writer` if easy, or use a temp file.*

## Out of Scope
- Client-side JS rendering (Chart.js, etc.).
- WebSocket or complex state management.
