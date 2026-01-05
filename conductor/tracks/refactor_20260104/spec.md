# Track Spec: Refactor code into a structured package and implement robust error handling

## Goal
Improve the codebase's maintainability, testability, and reliability by refactoring the monolithic `main.go` into a modular package structure and enhancing error handling.

## Requirements
- Create a `internal/gateway` package to house all T-Mobile Gateway interaction logic.
- Move existing data structures (`GatewayResponse`, `DeviceInfo`, etc.) and the `fetchStats` function to the new package.
- Enhance error handling:
    - Implement a retry mechanism for transient network errors.
    - Define custom error types for different failure scenarios (e.g., Timeout, UnmarshalError, APIError).
- Update `main.go` to import and use the `internal/gateway` package.
- Ensure the live dashboard and color-coding remain functional.
- Maintain or exceed 80% code coverage for the new package.
