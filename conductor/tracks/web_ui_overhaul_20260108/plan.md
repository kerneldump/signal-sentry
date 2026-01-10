# Track Plan: Web UI Time Range & Filtering Overhaul

## Phase 1: Backend Parsing Logic [checkpoint: 9656177]
- [x] Task: Create tests for `internal/web` to verify parsing of new query parameters (`range`, `start`, `end`). f6c8e04
- [x] Task: Update `internal/web/server.go` (or `request_handler.go`) to implement the parameter parsing logic. f6c8e04
- [x] Task: Conductor - User Manual Verification 'Backend Parsing Logic' (Protocol in workflow.md) f6c8e04

## Phase 2: Template Refactoring & Custom Input [checkpoint: 3ec3864]
- [x] Task: Update the Go templates (HTML) for the toolbar. cb2caab
- [x] Task: Conductor - User Manual Verification 'Template Refactoring & Custom Input' (Protocol in workflow.md) cb2caab

## Phase 3: Frontend Interaction (JavaScript) [checkpoint: 48f7369]
- [x] Task: Implement JavaScript logic for "Auto-Update". 30b92ee
- [x] Task: Implement "Mutual Exclusivity" logic. 30b92ee
- [x] Task: Conductor - User Manual Verification 'Frontend Interaction (JavaScript)' (Protocol in workflow.md) 30b92ee

## Phase 4: Integration & Polish [checkpoint: bf34b1d]
- [x] Task: Verify the entire flow.
- [x] Task: Polish the CSS (Bootstrap) to ensure the new inputs look good on mobile and desktop.
- [x] Task: Conductor - User Manual Verification 'Integration & Polish' (Protocol in workflow.md) bf34b1d
