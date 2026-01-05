# Track Spec: Add configurable refresh interval via command-line flag

## Goal
Allow users to customize how often the signal stats are refreshed by providing a command-line flag.

## Requirements
- Add a `-interval` command-line flag.
- The flag should accept an integer representing seconds.
- The default value must be `5` seconds.
- The tool must validate the input:
    - If the value is less than or equal to `0`, the tool must print an error message and exit with a non-zero status.
    - If the input is not a valid integer, the tool must exit with an error (standard `flag` package behavior).
- Use the standard Go `flag` package for implementation.
- The main polling loop must use the user-provided interval.

## Technical Details
- Use `flag.Int("interval", 5, "...")` or `flag.IntVar`.
- Validation should occur immediately after `flag.Parse()`.
- Error messages should be clear and technical, as per product guidelines.
