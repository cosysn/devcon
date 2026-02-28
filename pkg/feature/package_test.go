package feature

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageFeature(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	os.WriteFile(filepath.Join(dir, "devcontainer-feature.json"), []byte(`{"id": "test"}`), 0644)
	os.WriteFile(filepath.Join(dir, "install.sh"), []byte("#!/bin/bash\necho hello"), 0755)

	output := filepath.Join(t.TempDir(), "output.tar.gz")

	err := PackageFeature(dir, output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("output file should exist")
	}
}
