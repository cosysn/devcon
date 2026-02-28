package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DevcontainerConfig represents the devcontainer.json configuration
type DevcontainerConfig struct {
	Image         string                 `json:"image,omitempty"`
	Dockerfile    string                 `json:"dockerFile,omitempty"`
	Build         *BuildConfig           `json:"build,omitempty"`
	Features      map[string]interface{} `json:"features,omitempty"`
	Env           map[string]string      `json:"containerEnv,omitempty"`
	Mounts        []string               `json:"mounts,omitempty"`
	Ports         []int                  `json:"forwardPorts,omitempty"`
	RemoteUser    string                 `json:"remoteUser,omitempty"`
	Extends       string                 `json:"extends,omitempty"`
}

type BuildConfig struct {
	Dockerfile string `json:"dockerfile,omitempty"`
	Context    string `json:"context,omitempty"`
}

func ParseDevcontainer(dir string) (*DevcontainerConfig, error) {
	path := filepath.Join(dir, ".devcontainer", "devcontainer.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read devcontainer.json: %w", err)
	}

	parsed, err := ParseJSONC(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSONC: %w", err)
	}

	// Simple JSON unmarshal into struct
	jsonData, err := json.Marshal(parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parsed config: %w", err)
	}
	var config DevcontainerConfig
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
