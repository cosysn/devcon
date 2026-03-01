## Why

The current devcon CLI commands have minimal output - they only print basic status messages without progress indicators, detailed error context, or machine-readable formats. This makes it difficult for users to understand what's happening during long operations (like building images), debug issues, or integrate with other tools. Adding rich output capabilities will significantly improve the user experience.

## What Changes

- Add `--verbose` flag to all commands for detailed step-by-step output
- Add `--output json` flag to all commands for structured machine-readable output
- Enhance error messages with contextual information and actionable suggestions
- Add progress indicators for long-running operations (building, pulling images)
- Add `--quiet` flag to suppress non-essential output
- Create a unified output package to manage different output modes

## Capabilities

### New Capabilities
- `verbose-output`: Detailed step-by-step logging for all commands, showing each operation as it happens
- `structured-output`: JSON output mode for programmatic consumption and CI/CD integration
- `enhanced-errors`: Rich error messages with context, suggestions, and error codes
- `progress-indicators`: Visual progress feedback for long-running operations

### Modified Capabilities
<!-- No existing spec-level behavior changes -->

## Impact

- New package: `pkg/output` - unified output handling
- Modified commands: `up`, `build`, `run`, `config`, `features`, `inspect`
- Dependencies: May add a progress bar library (e.g., pb or similar)
