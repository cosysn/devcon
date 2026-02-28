package config

import (
	"encoding/json"
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
		return nil, err
	}

	parsed, err := ParseJSONC(string(data))
	if err != nil {
		return nil, err
	}

	// Simple JSON unmarshal into struct
	jsonData, _ := json.Marshal(parsed)
	var config DevcontainerConfig
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
