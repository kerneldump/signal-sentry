# Technology Stack - Signal Sentry

## Core Technologies
- **Programming Language:** Go (Golang)
- **Version:** 1.25.5 (as specified in `go.mod`)

## Key Libraries & Frameworks
- **Standard Library (`net/http`):** Used for all communication with the gateway's REST API.
- **Standard Library (`encoding/json`):** Handles parsing of the gateway's JSON responses into structured data.
- **Standard Library (`time`):** Manages the polling interval for real-time updates and timeouts.
- **Bootstrap (CSS Framework):** Used for styling the web dashboard (loaded via CDN).
- **Gonum/plot:** Used for generating high-resolution signal analysis charts.

## Architecture
- **Modular Structure:** Core gateway interaction logic is isolated in the `internal/gateway` package.
- **Resilient Client:** The gateway client implements an automated retry mechanism with exponential backoff for transient network issues.
- **Flexible Logging:** A pluggable logging system in `internal/logger` supports multiple output formats (JSON, CSV).

## Infrastructure & External APIs
- **T-Mobile Gateway API:** Interacts with the local API endpoint at `http://192.168.12.1/TMI/v1/gateway?get=all`.

## Development & Build Tools
- **Makefile:** Provides standardized commands for building, running, and testing the application.
- **Go Modules:** Manages project dependencies and versions.
