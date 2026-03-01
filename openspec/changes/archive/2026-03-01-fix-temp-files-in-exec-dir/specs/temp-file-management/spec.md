## ADDED Requirements

### Requirement: Temp files MUST be created in system temp directory

When the devcon CLI creates temporary files for feature packaging and publishing, those files MUST be created in the system's temp directory (e.g., `/tmp` on Linux, `%TEMP%` on Windows), not in the current working directory.

#### Scenario: PublishFeature creates temp file in system temp dir
- **WHEN** `feature.PublishFeature` is called to publish a feature
- **THEN** the temporary tarball file is created in the system temp directory (os.TempDir())
- **AND** the temp file has a unique name with "feature-" prefix

#### Scenario: Temp file is cleaned up after successful publish
- **WHEN** `feature.PublishFeature` completes successfully
- **THEN** the temporary tarball file is removed from the system temp directory

#### Scenario: Temp file is cleaned up after failed publish
- **WHEN** `feature.PublishFeature` fails with an error
- **THEN** the temporary tarball file is removed from the system temp directory
