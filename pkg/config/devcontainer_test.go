package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDevcontainer(t *testing.T) {
	dir := t.TempDir()
	devcontainerDir := filepath.Join(dir, ".devcontainer")
	os.MkdirAll(devcontainerDir, 0755)

	config := `{
        "image": "mcr.microsoft.com/devcontainers/base:ubuntu",
        "features": {
            "node": {}
        }
    }`

	os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(config), 0644)

	result, err := ParseDevcontainer(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Image != "mcr.microsoft.com/devcontainers/base:ubuntu" {
		t.Errorf("expected image, got %v", result.Image)
	}
}
