package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Feature struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name,omitempty"`
	Version     string                 `json:"version,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	DependsOn   []string               `json:"dependsOn,omitempty"`
	ContainerEnv map[string]string    `json:"containerEnv,omitempty"`
}

type FeatureDefinition struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Version   string                 `json:"version"`
	DependsOn []string               `json:"dependsOn"`
	Options   map[string]OptionSpec  `json:"options"`
}

type OptionSpec struct {
	Type        string      `json:"type"`
	Default     interface{} `json:"default"`
	Description string      `json:"description,omitempty"`
}

func ParseFeatureDefinition(dir string) (*FeatureDefinition, error) {
	path := filepath.Join(dir, "devcontainer-feature.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var def FeatureDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, err
	}

	return &def, nil
}
