# Track Spec: Add JSON and CSV file logging support

## Goal
Enable the application to log signal statistics to a file in either JSON or CSV format while continuing to display the real-time dashboard.

## Requirements
- **Flags:**
    - `-format`: Specifies the logging format. Accepted values: `json`, `csv`. Default: "" (logging disabled).
    - `-output`: Specifies the output filename. Default: `signal-data.json` (if format is json) or `signal-data.csv` (if format is csv).
- **Behavior:**
    - If `-format` is provided, the application must open the specified (or default) file for appending.
    - If the file does not exist, it should be created.
    - **Dashboard:** The CLI dashboard (stdout) must continue to function normally.
    - **JSON Format:** Each log entry should be a single line of valid JSON (NDJSON format) containing the full `GatewayResponse` data structure.
    - **CSV Format:**
        - Must include a header row if the file is new or empty.
        - Columns should include: Timestamp (ISO8601), 5G Band, 5G RSRP, 5G SINR, 5G Bars, 4G Band, 4G RSRP, 4G SINR, 4G Bars.
- **Error Handling:**
    - Exit with error if an invalid format is provided.
    - Exit with error if the file cannot be opened or written to.

## Technical Details
- Use `flag` package (integrate with existing flags).
- Create a new `logger` package or struct to handle file I/O and formatting.
- Ensure thread safety if strictly necessary (though this app is single-threaded main loop).
