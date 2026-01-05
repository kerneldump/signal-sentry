# Track Spec: Interactive Live View (CLI)

## Goal
Implement a rich, interactive Terminal User Interface (TUI) for `signal-sentry` using the `bubbletea` library. This mode will provide a live, non-scrolling view of the signal statistics with a history buffer and runtime controls.

## Key Features
1.  **Live View Mode:** Activated via a `-live` CLI flag (or potentially default if desired later, currently flag-gated).
2.  **Rolling Buffer:** Display a history of the last 30 signal stats.
    -   Newest stat at the top.
    -   Older stats push down.
3.  **Responsiveness:**
    -   If the terminal window height is insufficient (< 30 lines + header), truncate the older history lines to fit the screen.
4.  **Runtime Controls:**
    -   **Refresh Rate:** Users can dynamically adjust the signal polling interval using keys (e.g., `+`/`-` or `ArrowUp`/`ArrowDown`).
    -   **Exit:** Standard `q` or `Ctrl+C` to quit.
5.  **Tech Stack:**
    -   Library: `github.com/charmbracelet/bubbletea` (Go).
    -   Styling: Basic ANSI or `lipgloss` if needed for column alignment (though `bubbletea` alone is likely sufficient for structure).

## Integration
-   The UI will replace the standard `fmt.Println` loop when active.
-   It must integrate with the existing `gateway` package for data fetching.
-   It must respect the initial configuration (from flags/config).

## Constraints
-   Must work on macOS (Darwin) as per user context.
-   Must handle resize events gracefully.
