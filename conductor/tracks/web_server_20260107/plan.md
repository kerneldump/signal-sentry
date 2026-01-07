# Track Plan: Lightweight Web Server for Live Charts

## Phase 1: Subcommand & Scaffolding
- [x] Task: Create `cmd/web/main.go` (or logic within `main.go` to handle `web` subcommand). [Implementation complete]
- [x] Task: Define `runWeb` function in `main.go`. [Implementation complete]
- [x] Task: Create `internal/web` package to hold handler logic. [Implementation complete]

## Phase 2: Chart Generation Refactoring
- [x] Task: Refactor `internal/charting/charting.go`.
    -   Renamed `Generate` to `GenerateToWriter` (and kept wrapper).
    -   Modified to write to `io.Writer`.

## Phase 3: HTTP Handlers
- [x] Task: Implement `handleChart` in `internal/web`. [Implementation complete]
- [x] Task: Implement `handleIndex` in `internal/web`. [Implementation complete]

## Phase 4: Integration & Polish
- [x] Task: Wire up handlers in `main.go` and start `http.ListenAndServe`. [Implementation complete]
- [x] Task: Add auto-refresh (meta tag or JS) to the HTML. [Implementation complete]
- [x] Task: Verify: Run `signal-sentry web` and test in browser. [Verified via curl] [checkpoint: 06adec7]