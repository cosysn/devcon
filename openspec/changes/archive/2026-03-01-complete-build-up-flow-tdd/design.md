## Context

The devcon CLI has build and up commands that work but need:
1. Fixes to make all E2E tests pass (especially E2E-004 with extends fixture issue)
2. Ensure unit tests continue to pass
3. Use TDD approach for any new development

Current state:
- `devcon build` - builds devcontainer images
- `devcon up` - builds and starts containers
- Unit tests in `pkg/config` and `pkg/feature` pass
- E2E tests exist with 34 test cases, but E2E-004 fails

## Goals / Non-Goals

**Goals:**
- Fix failing E2E test E2E-004 (extends fixture path issue)
- Ensure all E2E tests pass (`make test-e2e`)
- Ensure all unit tests pass (`make test-unit`)
- Use TDD for any new feature development

**Non-Goals:**
- Major refactoring of existing working code
- Adding new commands (unless required by tests)
- Changing the architecture significantly

## Decisions

1. **Fix test fixtures first** - The E2E-004 failure is due to incorrect fixture path (`base.json` location doesn't match what `devcontainer.json` expects)
2. **Run tests before implementing** - Use TDD: write failing test, then fix/implement
3. **Keep existing architecture** - The build/up flow works, just needs test fixes

## Risks / Trade-offs

- **Risk**: Fixing test fixtures may expose real issues in code
  - **Mitigation**: Investigate root cause, fix code if needed

- **Risk**: E2E tests may be flaky with Docker
  - **Mitigation**: Use proper timeouts and cleanup in test scripts
