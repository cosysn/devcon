package config

// Feature represents a devcontainer feature configuration
type Feature struct {
    Version string `json:"version,omitempty"`
    Options map[string]string `json:"options,omitempty"`
}
