# Specification: Web UI Time Range & Filtering Overhaul

## Overview
Update the Signal Sentry web interface to provide a more streamlined set of preset time ranges and introduce advanced filtering capabilities through custom relative durations and absolute date ranges.

## User Interaction
- **Simplified Presets:** The toolbar will now feature a focused set of buttons: `10m`, `1h`, `6h`, `24h`, and `Max`.
- **Custom Relative Duration:** A text input field on the toolbar allows users to type specific durations (e.g., `2h`, `45m`, `90s`).
- **Date Range Selection:** Two native HTML `datetime-local` pickers will be provided for selecting a specific `Start` and `End` time window.
- **Mutual Exclusivity:** Selecting a preset or entering a custom relative duration will clear the absolute date range inputs. Conversely, selecting a date range will clear the relative duration selection.
- **Auto-Update:** The chart will automatically refresh when a preset is clicked, a date is changed, or the custom duration input is modified.

## Functional Requirements
- **Frontend (HTML/Templates):**
    - Refactor the button group in the web template to the new preset values.
    - Add the custom duration text input.
    - Add the `datetime-local` start/end inputs.
    - Implement JavaScript logic to handle the mutual exclusivity and auto-submission of the form (or AJAX refresh).
- **Backend (Go):**
    - Update `internal/web/server.go` to parse the new query parameters: `range` (relative), `start` (absolute), and `end` (absolute).
    - Map these parameters to the existing `analysis.TimeFilter` logic.
    - Ensure the `handleIndex` and `handleChart` functions correctly utilize these filters.

## Non-Functional Requirements
- **UI Aesthetic:** Maintain the current Bootstrap-based look and feel.
- **Responsiveness:** Ensure the expanded toolbar remains usable on smaller screens (mobile/tablet).

## Acceptance Criteria
- [ ] Preset buttons `10m`, `1h`, `6h`, `24h`, `Max` correctly filter data.
- [ ] Custom text input (e.g., `2h`) correctly filters data.
- [ ] Native date pickers correctly filter data between the selected timestamps.
- [ ] Inputs correctly reset/clear each other to maintain mutual exclusivity.
- [ ] The chart updates automatically upon selection change.

## Out of Scope
- Server-side validation of logical date order (e.g., ensuring Start < End) beyond standard parsing.
- Persistence of custom filter settings across browser sessions.
