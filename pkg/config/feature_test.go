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
