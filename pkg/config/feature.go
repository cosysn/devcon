package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
// It checks for local features and validates they exist.
func ResolveFeatures(dir string, features map[string]interface{}) error {
	if len(features) == 0 {
		return nil
	}

	// Check both possible locations for features:
	// 1. .devcontainer/features/<feature-id>/ (standard location)
	// 2. features/<feature-id>/ (alternative location)
	featuresDirs := []string{
		filepath.Join(dir, ".devcontainer", "features"),
		filepath.Join(dir, "features"),
	}

	for featureID := range features {
		// Check if it's an OCI reference (contains / or :)
		if isOCIReference(featureID) {
			// For OCI references, we can't validate locally - assume valid
			continue
		}

		// Check local features in both possible locations
		found := false
		for _, featuresDir := range featuresDirs {
			featurePath := filepath.Join(featuresDir, featureID)
			if _, err := os.Stat(featurePath); err == nil {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("feature %q not found in .devcontainer/features/ or features/", featureID)
		}
	}

	return nil
}

// isOCIReference checks if a string looks like an OCI registry reference
func isOCIReference(s string) bool {
	// OCI references typically contain / or : (for tag)
	// Examples: ghcr.io/user/feature, localhost:5000/feature:v1
	return strings.Contains(s, "/") || strings.Contains(s, ":")
}
