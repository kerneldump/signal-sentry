# Track Spec: Ping Reliability Investigation

## Problem
The user reports observing ~8 dropped packets in `stats.log` over 24+ hours, while other devices on the same network reported 0 loss. This suggests potential false positives in the `signal-sentry` ping implementation.

## Areas to Investigate

### 1. Timeout Logic
*   **Current:** `pinger.Timeout = p.Interval`.
*   **Hypothesis:** If the interval is 1s, and a ping takes 1.01s (or the Go scheduler is busy), the library might kill the ping context before the reply is processed, counting it as lost.
*   **Fix:** Decouple `Timeout` from `Interval`. A ping should probably have a dedicated timeout (e.g., 2s) independent of the polling ticker, or at least be slightly flexible.

### 2. Startup/Privilege Issues
*   **Current:** The code re-creates a `probing.NewPinger` *every single tick*.
*   **Hypothesis:** Re-initializing the raw socket every second might hit OS rate limits, file descriptor limits, or incur startup overhead that causes occasional drops.
*   **Fix:** Use a long-running Pinger instance that pings continuously (standard usage of `pro-bing`), rather than "One-Shot" pinging in a loop.

### 3. Privilege Drop
*   **Observation:** The code handles "permission denied" errors. If `sudo` privileges are momentarily dropped or confused (unlikely but possible with some OS policies), it might fail.

## Goal
Identify the root cause and propose a refactor to `internal/pinger`.

## Plan
1.  **Audit:** Review `internal/pinger/pinger.go` code.
2.  **Experiment:** Create a small reproduction script to stress-test the "One-Shot" vs "Long-Running" approach.
3.  **Propose Fix:** Likely switching to a continuous pinger or adjusting timeouts.
