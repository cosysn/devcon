## ADDED Requirements

### Requirement: Port attributes configuration
The devcontainer configuration SHALL support detailed port forwarding attributes.

#### Scenario: Port attributes specified
- **WHEN** devcontainer.json specifies portsAttributes
- **THEN** the configuration SHALL store port-specific attributes

#### Scenario: Other ports attributes specified
- **WHEN** devcontainer.json specifies otherPortsAttributes
- **THEN** the configuration SHALL store default attributes for ports not in portsAttributes

#### Scenario: Port label specified
- **WHEN** port attribute includes label property
- **THEN** the label SHALL be stored for display purposes

#### Scenario: Port protocol specified
- **WHEN** port attribute includes protocol property
- **THEN** the protocol SHALL be stored (tcp, http, https)

#### Scenario: Port auto-forward action specified
- **WHEN** port attribute includes onAutoForward property
- **THEN** the auto-forward action SHALL be stored (notify, openBrowser, openBrowserOnce, openPreview, silent, ignore)
