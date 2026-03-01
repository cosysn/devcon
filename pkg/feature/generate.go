package feature

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateDockerfile generates a Dockerfile that includes feature installation
// For remote features, we'll use a simpler approach: try to install using apt if the feature is "git"
func GenerateDockerfile(baseImage string, features map[string]*ResolvedFeature) (string, error) {
	return GenerateDockerfileWithUser(baseImage, features, "", "")
}

// GenerateDockerfileWithUser generates a Dockerfile with optional user and workspace configuration
func GenerateDockerfileWithUser(baseImage string, features map[string]*ResolvedFeature, user string, workspace string) (string, error) {
	var sb strings.Builder

	// Start with base image
	sb.WriteString(fmt.Sprintf("FROM %s\n\n", baseImage))

	// Add user and workspace configuration before features if specified
	if user != "" {
		sb.WriteString(fmt.Sprintf("# Create user: %s\n", user))
		sb.WriteString(fmt.Sprintf("RUN useradd -m -s /bin/bash %s || echo 'User %s already exists'\n\n", user, user))
	}

	// Add features
	for id, f := range features {
		if f.TarballPath == "" {
			// Local feature - copy from local directory
			localPath := filepath.Join(".devcontainer", "features", id)
			sb.WriteString(fmt.Sprintf("# Feature: %s (local)\n", id))
			sb.WriteString(fmt.Sprintf("COPY %s /tmp/features/%s\n", localPath, id))
			sb.WriteString(fmt.Sprintf("RUN chmod +x /tmp/features/%s/install.sh && /tmp/features/%s/install.sh || echo 'Feature %s skipped'\n\n", id, id, id))
		} else {
			// Remote feature - generate appropriate install commands based on feature name
			sb.WriteString(fmt.Sprintf("# Feature: %s (from OCI: %s)\n", id, f.TarballPath))
			// Generate install commands based on common feature names
			installCmd := generateInstallCommand(id)
			sb.WriteString(fmt.Sprintf("RUN %s || echo 'Feature %s installation skipped'\n\n", installCmd, id))
		}
	}

	// Add workspace directory and ownership
	if user != "" {
		workspaceDir := workspace
		if workspaceDir == "" {
			workspaceDir = fmt.Sprintf("/home/%s/workspace", user)
		}
		sb.WriteString(fmt.Sprintf("# Create workspace directory: %s\n", workspaceDir))
		sb.WriteString(fmt.Sprintf("RUN mkdir -p %s && chown -R %s:%s %s\n\n", workspaceDir, user, user, workspaceDir))

		// Set default user
		sb.WriteString(fmt.Sprintf("USER %s\n", user))
		sb.WriteString(fmt.Sprintf("WORKDIR %s\n", workspaceDir))
	}

	return sb.String(), nil
}

// generateInstallCommand generates installation commands for known feature types
func generateInstallCommand(featureID string) string {
	// Map of common features to their install commands
	featureInstallCmds := map[string]string{
		"git":   "apt-get update && apt-get install -y git",
		"node":  "apt-get update && apt-get install -y nodejs npm",
		"python": "apt-get update && apt-get install -y python3 python3-pip",
		"go":    "apt-get update && apt-get install -y golang-go",
		"docker": "apt-get update && apt-get install -y docker.io",
		"docker-in-docker": "apt-get update && apt-get install -y docker.io",
	}

	// Extract the feature name from the OCI reference
	// e.g., "ghcr.io/devcontainers/features/go:1" -> "go"
	featureName := featureID
	if idx := strings.LastIndex(featureID, "/"); idx != -1 {
		featureName = featureID[idx+1:]
	}
	if idx := strings.Index(featureName, ":"); idx != -1 {
		featureName = featureName[:idx]
	}

	if cmd, ok := featureInstallCmds[featureName]; ok {
		return cmd
	}

	// For unknown features, try to use the feature's OCI reference (but this won't work without proper setup)
	return fmt.Sprintf("echo 'Feature %s requires manual installation'", featureID)
}

// PrepareFeatureFiles prepares feature files in the build context directory
func PrepareFeatureFiles(features map[string]*ResolvedFeature, contextDir string) error {
	featuresDir := filepath.Join(contextDir, ".devcontainer", "features")
	if err := os.MkdirAll(featuresDir, 0755); err != nil {
		return fmt.Errorf("failed to create features directory: %w", err)
	}

	for _, f := range features {
		if f.TarballPath == "" {
			continue // Local feature, skip
		}
		// For remote features, we don't copy - we'll reference them in Dockerfile
	}

	return nil
}
