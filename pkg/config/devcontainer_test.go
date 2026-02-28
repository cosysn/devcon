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
