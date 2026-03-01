package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

// TopologicalSort returns the features in dependency order using Kahn's algorithm
func TopologicalSort(features map[string]*FeatureDefinition) ([]string, error) {
	// Build dependency graph
	deps := make(map[string]map[string]bool)
	allDeps := make(map[string][]string)

	for id, f := range features {
		deps[id] = make(map[string]bool)
		for _, d := range f.DependsOn {
			deps[id][d] = true
		}
		allDeps[id] = f.DependsOn
	}

	// Calculate in-degrees (number of dependencies for each feature)
	inDegree := make(map[string]int)
	for id := range features {
		inDegree[id] = len(allDeps[id])
	}

	// Find all nodes with no incoming edges
	var queue []string
	for id, d := range inDegree {
		if d == 0 {
			queue = append(queue, id)
		}
	}

	// Process queue
	var result []string
	for len(queue) > 0 {
		sort.Strings(queue)
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		for id, depSet := range deps {
			if depSet[current] {
				inDegree[id]--
				if inDegree[id] == 0 {
					queue = append(queue, id)
				}
			}
		}
	}

	if len(result) != len(features) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return result, nil
}

// ResolveFeatures validates and resolves features from the devcontainer config.
// It checks for local features. Remote features (shorthand names or OCI references)
// are not validated here - they'll be resolved during build.
func ResolveFeatures(dir string, features map[string]interface{}) error {
	if len(features) == 0 {
		return nil
	}

	// Check both possible locations for local features:
	// 1. .devcontainer/features/<feature-id>/ (standard location)
	// 2. features/<feature-id>/ (alternative location)
	featuresDirs := []string{
		filepath.Join(dir, ".devcontainer", "features"),
		filepath.Join(dir, "features"),
	}

	for featureID := range features {
		// Check local features in both possible locations
		found := false
		for _, featuresDir := range featuresDirs {
			featurePath := filepath.Join(featuresDir, featureID)
			if _, err := os.Stat(featurePath); err == nil {
				found = true
				break
			}
		}

		// If not found locally, assume it's a remote feature (shorthand name or OCI reference)
		// We don't fail here because:
		// - Shorthand names like "docker-in-docker" should be resolved from remote registry
		// - Full OCI references like "ghcr.io/user/feature" are remote
		// The actual resolution happens during the build process
		if !found {
			// Just log/continue - remote features will be resolved during build
			continue
		}
	}

	return nil
}
