## Why

When running devcon commands like `devcon build`, temporary files are being created in the current execution directory (current working directory) instead of being properly managed in system temp directories. This pollutes the user's working directory with temporary files that should be cleaned up automatically.

## What Changes

- Fix `PublishFeature` in `pkg/feature/package.go` to use proper temp directory handling
- Ensure all temp files are created in system temp directories (e.g., `/tmp` on Linux)
- Add proper cleanup for temp files to prevent orphaned files
- Review and fix any other locations where files might be created in the execution directory

## Capabilities

### New Capabilities
- `temp-file-management`: Ensure all temp files are created in system temp directories with proper cleanup

### Modified Capabilities
- None

## Impact

- `pkg/feature/package.go`: Fix temp file creation in `PublishFeature` function
- Any other files that may create temp files in execution directory
