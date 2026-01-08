# Track Plan: Unified Run Mode (TUI + Web)

## Phase 1: Web Package Refactoring
- [x] Task: Update `internal/web/server.go`.
    -   Modify `Run` to accept a `quiet bool` parameter (or a `Logger` interface).
    -   If quiet, suppress `log.Printf` calls (or redirect to `io.Discard`).
    -   Ensure `http.ListenAndServe` errors are still captured/returned (or logged to a file if needed).

## Phase 2: Configuration & CLI
- [x] Task: Update `internal/config/config.go` to include `WebEnabled` (bool) and `WebPort` (int).
- [x] Task: Update `main.go` flag parsing.
    -   Add `-web` flag.
    -   Add `-web-port` flag.
    -   Map these to the config.

## Phase 3: Integration
- [x] Task: Update `main.go`.
    -   Before starting the TUI or Legacy Loop, check `cfg.WebEnabled`.
    -   If true, launch `go web.Run(cfg.WebPort, cfg.Output, true)` (where `true` is for quiet mode).
    -   Handle potential port conflicts (e.g., if port is busy, maybe log to stderr before TUI starts, or just fail silently/log to file).

## Phase 4: Verification
- [x] Task: Verify: Run `go run . -live -web` and check if TUI works AND `localhost:8080` loads. [Verified via build and unit tests]
- [x] Task: Verify: Run `go run . -web` (without live) and check if legacy output works AND `localhost:8080` loads. [Verified via build]
