# Technology Stack - Signal Sentry

## Core Technologies
- **Programming Language:** Go (Golang)
- **Version:** 1.25.5 (as specified in `go.mod`)

## Key Libraries & Frameworks
- **Standard Library (`net/http`):** Used for all communication with the gateway's REST API.
- **Standard Library (`encoding/json`):** Handles parsing of the gateway's JSON responses into structured data.
- **Standard Library (`time`):** Manages the polling interval for real-time updates and timeouts.

## Infrastructure & External APIs
- **T-Mobile Gateway API:** Interacts with the local API endpoint at `http://192.168.12.1/TMI/v1/gateway?get=all`.

## Development & Build Tools
- **Makefile:** Provides standardized commands for building, running, and testing the application.
- **Go Modules:** Manages project dependencies and versions.
