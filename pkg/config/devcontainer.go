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
	Image              string                 `json:"image,omitempty"`
	Dockerfile         string                 `json:"dockerFile,omitempty"`
	Build             *BuildConfig           `json:"build,omitempty"`
	Features          map[string]interface{} `json:"features,omitempty"`
	Env               map[string]string      `json:"containerEnv,omitempty"`
	Mounts            []string               `json:"mounts,omitempty"`
	Ports             []int                  `json:"forwardPorts,omitempty"`
	RemoteUser        string                 `json:"remoteUser,omitempty"`
	User              string                 `json:"user,omitempty"`
	Workspace         string                 `json:"workspace,omitempty"`
	Extends           string                 `json:"extends,omitempty"`
	OnCreateCommand   string                 `json:"onCreateCommand,omitempty"`
	PostCreateCommand string                 `json:"postCreateCommand,omitempty"`
	PostStartCommand  string                 `json:"postStartCommand,omitempty"`
	Customizations    *Customizations        `json:"customizations,omitempty"`
}

type BuildConfig struct {
	Dockerfile string `json:"dockerfile,omitempty"`
	Context    string `json:"context,omitempty"`
}

type Customizations struct {
	VSCode *VSCodeCustomization `json:"vscode,omitempty"`
}

type VSCodeCustomization struct {
	Extensions []string               `json:"extensions,omitempty"`
	Settings   map[string]interface{} `json:"settings,omitempty"`
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

	// If Dockerfile is specified in build config, promote it to top-level
	if config.Dockerfile == "" && config.Build != nil && config.Build.Dockerfile != "" {
		config.Dockerfile = config.Build.Dockerfile
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

	// Clean and verify path stays within .devcontainer
	cleanPath := filepath.Clean(fullPath)
	devcontainerDir := filepath.Clean(filepath.Join(dir, ".devcontainer"))
	if !strings.HasPrefix(cleanPath, devcontainerDir) {
		return nil, fmt.Errorf("extends path escapes devcontainer directory: %s", fullPath)
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
		return nil, fmt.Errorf("failed to resolve base extends: %w", err)
	}
	base = *resolvedBase

	// If Dockerfile is specified in build config, promote it to top-level
	if base.Dockerfile == "" && base.Build != nil && base.Build.Dockerfile != "" {
		base.Dockerfile = base.Build.Dockerfile
	}

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
	if config.User != "" {
		base.User = config.User
	}
	if config.Workspace != "" {
		base.Workspace = config.Workspace
	}
	if config.Extends != "" {
		base.Extends = config.Extends
	}

	// Merge customizations
	if config.Customizations != nil {
		if base.Customizations == nil {
			base.Customizations = &Customizations{}
		}
		if config.Customizations.VSCode != nil {
			if base.Customizations.VSCode == nil {
				base.Customizations.VSCode = &VSCodeCustomization{}
			}
			// Merge extensions
			if len(config.Customizations.VSCode.Extensions) > 0 {
				base.Customizations.VSCode.Extensions = append(base.Customizations.VSCode.Extensions, config.Customizations.VSCode.Extensions...)
			}
			// Merge settings
			if config.Customizations.VSCode.Settings != nil {
				if base.Customizations.VSCode.Settings == nil {
					base.Customizations.VSCode.Settings = make(map[string]interface{})
				}
				for k, v := range config.Customizations.VSCode.Settings {
					base.Customizations.VSCode.Settings[k] = v
				}
			}
		}
	}

	return &base, nil
}
