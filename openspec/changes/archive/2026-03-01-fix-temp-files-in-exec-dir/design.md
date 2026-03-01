## Context

Currently, when running devcon commands (particularly `devcon build` and `devcon features publish`), temporary files may be created in the current execution directory instead of being properly managed in system temp directories. The issue is in `pkg/feature/package.go` where:

1. `PublishFeature` function uses `os.CreateTemp("", "feature-*.tar.gz")` - this should use system temp dir but needs verification
2. The code has proper defer cleanup with `os.Remove(tmpPath)`, but there may be edge cases where cleanup fails

## Goals / Non-Goals

**Goals:**
- Ensure all temp files are created in system temp directories (via `os.CreateTemp` with empty first arg or explicit temp dir)
- Verify proper cleanup of temp files even on error paths
- Prevent any temp files from being left in the execution directory

**Non-Goals:**
- Not changing the file output behavior (e.g., where `Dockerfile.with-features` is written - that's intentional build output)
- Not modifying the test file behavior (tests use `t.TempDir()` which is correct)

## Decisions

1. **Use system temp directory** - Confirm that `os.CreateTemp("", pattern)` uses the system's default temp directory correctly
2. **Explicit temp directory path** - Consider using `os.MkdirTemp(os.TempDir(), prefix)` for clarity and explicit control
3. **Enhanced error handling** - Ensure temp files are properly cleaned up even when errors occur during the publish process

## Risks / Trade-offs

- **Risk**: If the system's temp directory is not writable, the feature publish will fail
  - **Mitigation**: This is acceptable behavior - if temp is not writable, we should fail loudly rather than fall back to current directory
- **Risk**: Potential race condition if multiple publishes run simultaneously
  - **Mitigation**: `os.CreateTemp` handles this with unique filenames
