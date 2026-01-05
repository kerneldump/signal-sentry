# Track Spec: Always-on JSON Logging

## Goal
Ensure that `signal-sentry` automatically logs all signal statistics to a `stats.log` file in JSON format by default, regardless of other CLI flags. This creates a persistent history for future analysis.

## Requirements
1.  **Default Logging:**
    -   On startup, the application must open (or create) `stats.log` in append mode.
    -   Every signal fetch cycle must write the JSON representation of the data to this file.
    -   This happens *in addition* to any user-requested `-output` (CSV/JSON).
2.  **File Management:**
    -   Filename: `stats.log` (in the current working directory).
    -   Permissions: Standard read/write.
    -   Git Ignore: Add `stats.log` to `.gitignore`.
3.  **Concurrency:**
    -   Ensure thread safety if multiple loggers (user-requested + default) are writing (though they are likely writing to different files, so just standard file I/O safety).

## Integration
-   Modify `main.go` to initialize this "background logger" alongside the "user logger".
-   If the user *explicitly* asks to output to `stats.log` via flags, handle the collision gracefully (e.g., just use the one logger).
