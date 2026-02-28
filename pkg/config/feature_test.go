package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFeatureDefinition(t *testing.T) {
	dir := t.TempDir()

	config := `{
        "id": "node",
        "name": "Node.js",
        "version": "1.0.0",
        "dependsOn": ["git"],
        "options": {
            "version": {
                "type": "string",
                "default": "20"
            }
        }
    }`

	os.WriteFile(filepath.Join(dir, "devcontainer-feature.json"), []byte(config), 0644)

	result, err := ParseFeatureDefinition(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "node" {
		t.Errorf("expected id 'node', got %v", result.ID)
	}

	if len(result.DependsOn) != 1 || result.DependsOn[0] != "git" {
		t.Errorf("expected dependsOn [git], got %v", result.DependsOn)
	}
}

func TestParseFeatureDefinitionNotFound(t *testing.T) {
	_, err := ParseFeatureDefinition("/nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestTopologicalSort(t *testing.T) {
	tests := []struct {
		name      string
		features  map[string]*FeatureDefinition
		wantOrder []string
		wantErr   bool
	}{
		{
			name: "no dependencies",
			features: map[string]*FeatureDefinition{
				"a": {ID: "a", DependsOn: []string{}},
				"b": {ID: "b", DependsOn: []string{}},
			},
			wantErr: false,
		},
		{
			name: "linear dependencies",
			features: map[string]*FeatureDefinition{
				"a": {ID: "a", DependsOn: []string{}},
				"b": {ID: "b", DependsOn: []string{"a"}},
				"c": {ID: "c", DependsOn: []string{"b"}},
			},
			wantErr: false,
		},
		{
			name: "parallel dependencies",
			features: map[string]*FeatureDefinition{
				"a": {ID: "a", DependsOn: []string{}},
				"b": {ID: "b", DependsOn: []string{}},
				"c": {ID: "c", DependsOn: []string{"a", "b"}},
			},
			wantErr: false,
		},
		{
			name: "circular dependency",
			features: map[string]*FeatureDefinition{
				"a": {ID: "a", DependsOn: []string{"c"}},
				"b": {ID: "b", DependsOn: []string{"a"}},
				"c": {ID: "c", DependsOn: []string{"b"}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := TopologicalSort(tt.features)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantErr && len(order) != len(tt.features) {
				t.Errorf("got %d features, want %d", len(order), len(tt.features))
			}
		})
	}
}
