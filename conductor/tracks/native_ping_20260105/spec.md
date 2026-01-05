# Track Spec: Native Ping Implementation

## Goal
Replace the current `os/exec` "shell out" ping implementation with a native Go library solution. This improves portability, control, and efficiency, though it may require elevated privileges (`sudo`) to open raw sockets.

## Requirements
1.  **Library:** Use `github.com/prometheus-community/pro-bing` (formerly `go-ping/ping`), which is the industry standard for Go pingers.
2.  **Functionality:**
    -   Maintain the current behavior: Send 1 ping every second.
    -   Maintain statistics: Min, Avg, Max, StdDev, Packet Loss.
    -   Maintain support for `LastRTT` (the most recent ping value).
3.  **Privileges:**
    -   Acknowledge that `sudo` might be required on some systems (Linux/macOS) for ICMP.
    -   On macOS, it might work without sudo if unprivileged ICMP is enabled, but we will assume sudo is acceptable as per user request.
4.  **Integration:**
    -   Replace the logic in `internal/pinger/pinger.go`.
    -   Ensure `main.go` and `internal/ui` remain unaffected (interface should match).

## Implementation Details
-   The current `Pinger` struct runs a loop. We will modify `ping()` to use the library's `pinger.Count = 1` and `pinger.Run()`.
-   We need to translate the library's `Statistics` struct into our `models.PingStats`.
