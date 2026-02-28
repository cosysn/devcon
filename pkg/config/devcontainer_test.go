package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDevcontainer(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr bool
	}{
		{
			name: "basic image",
			config: `{"image": "mcr.microsoft.com/devcontainers/base:ubuntu"}`,
			wantErr: false,
		},
		{
			name: "with features",
			config: `{"image": "ubuntu", "features": {"node": {}}}`,
			wantErr: false,
		},
		{
			name: "with env",
			config: `{"image": "ubuntu", "containerEnv": {"VAR": "value"}}`,
			wantErr: false,
		},
		{
			name: "invalid json",
			config: `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			devcontainerDir := filepath.Join(dir, ".devcontainer")
			os.MkdirAll(devcontainerDir, 0755)
			os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(tt.config), 0644)

			_, err := ParseDevcontainer(dir)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestParseDevcontainerNotFound(t *testing.T) {
	_, err := ParseDevcontainer("/nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestResolveExtends(t *testing.T) {
	dir := t.TempDir()
	devcontainerDir := filepath.Join(dir, ".devcontainer")
	os.MkdirAll(devcontainerDir, 0755)

	// Create base config
	baseConfig := `{
        "image": "mcr.microsoft.com/devcontainers/base:ubuntu",
        "features": {
            "git": {}
        },
        "containerEnv": {
            "BASE": "value"
        }
    }`
	os.WriteFile(filepath.Join(devcontainerDir, "base.json"), []byte(baseConfig), 0644)

	// Create extending config
	config := `{
        "extends": "./base.json",
        "image": "custom:latest",
        "features": {
            "node": {}
        }
    }`
	os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(config), 0644)

	result, err := ParseDevcontainer(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err = ResolveExtends(dir, result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have custom image (config takes precedence)
	if result.Image != "custom:latest" {
		t.Errorf("expected image 'custom:latest', got %v", result.Image)
	}

	// Should have both git and node features
	if result.Features == nil {
		t.Fatal("features should not be nil")
	}
	if _, ok := result.Features["git"]; !ok {
		t.Error("features should include git")
	}
	if _, ok := result.Features["node"]; !ok {
		t.Error("features should include node")
	}

	// Should have merged env
	if result.Env == nil {
		t.Fatal("env should not be nil")
	}
	if v, ok := result.Env["BASE"]; !ok || v != "value" {
		t.Errorf("env should contain BASE=value, got %v", result.Env)
	}
}

func TestResolveExtendsNoExtends(t *testing.T) {
	dir := t.TempDir()
	devcontainerDir := filepath.Join(dir, ".devcontainer")
	os.MkdirAll(devcontainerDir, 0755)

	config := `{"image": "test:latest"}`
	os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(config), 0644)

	result, err := ParseDevcontainer(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err = ResolveExtends(dir, result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Image != "test:latest" {
		t.Errorf("expected image 'test:latest', got %v", result.Image)
	}
}

func TestResolveExtendsPathTraversal(t *testing.T) {
	dir := t.TempDir()
	devcontainerDir := filepath.Join(dir, ".devcontainer")
	os.MkdirAll(devcontainerDir, 0755)

	// Create config that tries to escape
	config := `{
        "extends": "../../../etc/passwd"
    }`
	os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(config), 0644)

	result, err := ParseDevcontainer(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = ResolveExtends(dir, result)
	if err == nil {
		t.Error("expected error for path traversal, got nil")
	}
}

func TestResolveExtendsMultipleLevels(t *testing.T) {
	dir := t.TempDir()
	devcontainerDir := filepath.Join(dir, ".devcontainer")
	os.MkdirAll(devcontainerDir, 0755)

	// Create level 1 config
	level1Config := `{
        "image": "level1:latest",
        "containerEnv": {"LEVEL": "1"}
    }`
	os.WriteFile(filepath.Join(devcontainerDir, "level1.json"), []byte(level1Config), 0644)

	// Create level 2 config extending level1
	level2Config := `{
        "extends": "./level1.json",
        "containerEnv": {"LEVEL": "2"}
    }`
	os.WriteFile(filepath.Join(devcontainerDir, "level2.json"), []byte(level2Config), 0644)

	// Create main config extending level2
	mainConfig := `{
        "extends": "./level2.json",
        "image": "main:latest"
    }`
	os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(mainConfig), 0644)

	result, err := ParseDevcontainer(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err = ResolveExtends(dir, result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have main image
	if result.Image != "main:latest" {
		t.Errorf("expected image 'main:latest', got %v", result.Image)
	}

	// Should have merged env (level2 should override level1)
	if result.Env == nil {
		t.Fatal("env should not be nil")
	}
	if result.Env["LEVEL"] != "2" {
		t.Errorf("expected LEVEL=2, got %v", result.Env["LEVEL"])
	}
}

func TestResolveExtendsNestedPath(t *testing.T) {
	dir := t.TempDir()
	devcontainerDir := filepath.Join(dir, ".devcontainer")
	os.MkdirAll(devcontainerDir, 0755)

	// Create subdirectory with base config
	subDir := filepath.Join(devcontainerDir, "nested")
	os.MkdirAll(subDir, 0755)
	baseConfig := `{"image": "nested-base:latest"}`
	os.WriteFile(filepath.Join(subDir, "base.json"), []byte(baseConfig), 0644)

	// Create main config extending nested path
	config := `{
        "extends": "./nested/base.json",
        "image": "custom:latest"
    }`
	os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(config), 0644)

	result, err := ParseDevcontainer(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err = ResolveExtends(dir, result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Image != "custom:latest" {
		t.Errorf("expected image 'custom:latest', got %v", result.Image)
	}
}
