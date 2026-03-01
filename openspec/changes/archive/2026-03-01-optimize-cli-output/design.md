## Context

The devcon CLI currently provides minimal output - basic Println statements without any verbose mode, structured output, or progress indicators. This affects all commands: `up`, `build`, `run`, `config`, `features`, and `inspect`.

Current issues:
- No verbose mode to see detailed step-by-step progress
- No JSON output for CI/CD integration
- Errors lack context and actionable suggestions
- Long operations (building images) show no progress feedback

## Goals / Non-Goals

**Goals:**
- Add `--verbose` flag for detailed logging
- Add `--output json` flag for structured machine-readable output
- Enhance error messages with context and suggestions
- Add progress indicators for long-running operations
- Create a unified `pkg/output` package

**Non-Goals:**
- Changing command behavior or functionality
- Adding new commands
- Modifying the underlying builder or container logic

## Decisions

### 1. Output Package Architecture
**Decision**: Create a new `pkg/output` package with a unified interface

**Rationale**: Centralizes all output logic, making it easy to add new output modes and maintain consistency across commands. Using an interface allows for easy testing.

**Alternatives Considered**:
- Directly modifying each command (rejected: code duplication, inconsistent behavior)
- Using a logging library (rejected: overkill for CLI, adds dependency)

### 2. Progress Indicator Implementation
**Decision**: Use a simple spinner/progress implementation with fallback for non-TTY

**Rationale**: Avoids external dependencies while providing visual feedback. Must handle non-interactive terminals gracefully.

**Alternatives Considered**:
- Using `pb` library (rejected: adds dependency)
- Using `spinner` library (rejected: adds dependency)

### 3. JSON Output Structure
**Decision**: Use a consistent envelope format for all JSON output

```json
{
  "success": true,
  "command": "build",
  "message": "...",
  "data": { ... }
}
```

**Rationale**: Consistent structure makes it easy to parse and handle errors programmatically.

## Risks / Trade-offs

- **Risk**: Progress indicators may not work in all environments (CI, piped output)
  - **Mitigation**: Detect TTY and fall back to simple text output

- **Risk**: JSON output may break existing scripts
  - **Mitigation**: JSON output only enabled with explicit `--output json` flag, not default

- **Risk**: Verbose output may be too verbose
  - **Mitigation**: Use appropriate log levels (info vs debug)

## Migration Plan

1. Add `pkg/output` package with Output interface
2. Update each command to use the output package
3. Add flags to root command (--verbose, --output, --quiet)
4. Commands inherit flags from root
5. Test each command with new output modes
