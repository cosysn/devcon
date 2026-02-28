package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// ResolveExtends resolves the extends property in devcontainer config
// It recursively merges the base config with the current config
func ResolveExtends(dir string, config *DevcontainerConfig) (*DevcontainerConfig, error) {
	if config.Extends == "" {
		return config, nil
	}

	// Resolve extends path - could be "./base.json" or just "base"
	extendsPath := config.Extends
	if !filepath.IsAbs(extendsPath) && !strings.HasPrefix(extendsPath, "./") && !strings.HasPrefix(extendsPath, "../") {
		extendsPath = "./" + extendsPath
	}

	fullPath := filepath.Join(dir, ".devcontainer", extendsPath)
	// Add .json if not present
	if !strings.HasSuffix(fullPath, ".json") {
		fullPath += ".json"
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read extended config %s: %w", fullPath, err)
	}

	parsed, err := ParseJSONC(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse extended config: %w", err)
	}

	var base DevcontainerConfig
	jsonData, err := json.Marshal(parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal base config: %w", err)
	}
	if err := json.Unmarshal(jsonData, &base); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base config: %w", err)
	}

	// Recursively resolve extends
	resolvedBase, err := ResolveExtends(dir, &base)
	if err != nil {
		return nil, err
	}
	base = *resolvedBase

	// Merge: base -> config (config takes precedence)
	if config.Image != "" {
		base.Image = config.Image
	}

	// Merge features
	if config.Features != nil {
		if base.Features == nil {
			base.Features = make(map[string]interface{})
		}
		for k, v := range config.Features {
			base.Features[k] = v
		}
	}

	// Merge env
	if config.Env != nil {
		if base.Env == nil {
			base.Env = make(map[string]string)
		}
		for k, v := range config.Env {
			base.Env[k] = v
		}
	}

	// Merge mounts
	if config.Mounts != nil {
		base.Mounts = append(base.Mounts, config.Mounts...)
	}

	// Merge ports
	if config.Ports != nil {
		base.Ports = append(base.Ports, config.Ports...)
	}

	// Merge other fields
	if config.Dockerfile != "" {
		base.Dockerfile = config.Dockerfile
	}
	if config.RemoteUser != "" {
		base.RemoteUser = config.RemoteUser
	}
	if config.Extends != "" {
		base.Extends = config.Extends
	}

	return &base, nil
}
