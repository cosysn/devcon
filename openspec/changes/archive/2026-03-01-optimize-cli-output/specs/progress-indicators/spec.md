## ADDED Requirements

### Requirement: Progress indicator during image build
The CLI SHALL display a progress indicator while building the container image.

#### Scenario: Build shows progress
- **WHEN** user runs `devcon build .`
- **THEN** a progress indicator (spinner or progress bar) is displayed during the build

#### Scenario: Build shows completion
- **WHEN** build completes successfully
- **THEN** progress indicator shows completion state and final image ID

### Requirement: Progress indicator during feature resolution
The CLI SHALL display progress while resolving and downloading features.

#### Scenario: Feature resolution shows progress
- **WHEN** features are being resolved
- **THEN** a progress indicator shows the resolution process

### Requirement: Graceful fallback for non-TTY
The CLI SHALL detect non-TTY environments and fall back to simple text output without progress indicators.

#### Scenario: Non-interactive environment
- **WHEN** output is piped or redirected (non-TTY)
- **THEN** plain text output is used instead of progress indicators
