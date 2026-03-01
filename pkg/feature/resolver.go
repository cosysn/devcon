package feature

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devcon/cli/pkg/config"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type Resolver struct {
	registryOpts []remote.Option
}

func NewResolver(registryOpts ...remote.Option) *Resolver {
	return &Resolver{
		registryOpts: registryOpts,
	}
}

// ResolveLocalFeatures reads Feature definitions from local .devcontainer/features directory
// It checks both .devcontainer/features/ and features/ directories
func (r *Resolver) ResolveLocalFeatures(dir string) (map[string]*config.FeatureDefinition, error) {
	result := make(map[string]*config.FeatureDefinition)

	// Check both possible locations
	featuresDirs := []string{
		filepath.Join(dir, ".devcontainer", "features"),
		filepath.Join(dir, "features"),
	}

	for _, featuresDir := range featuresDirs {
		entries, err := os.ReadDir(featuresDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			featurePath := filepath.Join(featuresDir, entry.Name())
			def, err := config.ParseFeatureDefinition(featurePath)
			if err != nil {
				return nil, fmt.Errorf("failed to parse feature %s: %w", entry.Name(), err)
			}

			// Avoid overwriting if duplicate
			if _, exists := result[def.ID]; !exists {
				result[def.ID] = def
			}
		}
	}

	return result, nil
}

// ResolveOCIFeature pulls and parses a Feature from OCI registry
func (r *Resolver) ResolveOCIFeature(ctx context.Context, ref string) (*config.FeatureDefinition, error) {
	refParsed, err := name.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	img, err := remote.Image(refParsed, r.registryOpts...)
	if err != nil {
		return nil, err
	}

	// Get config file
	cfg, err := img.ConfigFile()
	if err != nil {
		return nil, err
	}

	// Parse devcontainer-feature.json from label or annotation
	label := cfg.Config.Labels["devcontainer-feature"]
	if label == "" {
		return nil, fmt.Errorf("no devcontainer-feature label found in image")
	}

	var def config.FeatureDefinition
	if err := json.Unmarshal([]byte(label), &def); err != nil {
		return nil, fmt.Errorf("failed to parse feature definition: %w", err)
	}

	return &def, nil
}

// DefaultFeatureRegistry is the default registry for shorthand feature names
const DefaultFeatureRegistry = "ghcr.io/devcontainers/features"

// ConvertShorthandToOCI converts a shorthand feature name to full OCI reference
// e.g., "git" -> "ghcr.io/devcontainers/features/git:latest"
func ConvertShorthandToOCI(featureID string) string {
	// If it's already a full reference (contains /), return as-is
	if strings.Contains(featureID, "/") {
		return featureID
	}
	return fmt.Sprintf("%s/%s:latest", DefaultFeatureRegistry, featureID)
}

// IsShorthand returns true if the feature ID is a shorthand (no /)
func IsShorthand(featureID string) bool {
	return !strings.Contains(featureID, "/")
}

// ResolvedFeature represents a feature that has been resolved (downloaded)
type ResolvedFeature struct {
	Definition *config.FeatureDefinition
	TarballPath string
	InstallPath string
}

// ResolveAndDownload resolves all features from devcontainer config and downloads them
func (r *Resolver) ResolveAndDownload(ctx context.Context, dir string, features map[string]interface{}) (map[string]*ResolvedFeature, error) {
	if len(features) == 0 {
		return make(map[string]*ResolvedFeature), nil
	}

	result := make(map[string]*ResolvedFeature)

	// First, resolve local features
	localFeatures, err := r.ResolveLocalFeatures(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local features: %w", err)
	}

	// Add local features to result
	for id, def := range localFeatures {
		result[id] = &ResolvedFeature{
			Definition: def,
			TarballPath: "",
			InstallPath: filepath.Join(dir, ".devcontainer", "features", id),
		}
	}

	// Then resolve remote features (shorthand or full OCI references)
	for featureID := range features {
		// Skip if already resolved as local
		if _, ok := result[featureID]; ok {
			continue
		}

		// Convert shorthand to full OCI reference
		ociRef := ConvertShorthandToOCI(featureID)

		// Download the feature
		tarballPath, err := r.DownloadFeature(ctx, ociRef)
		if err != nil {
			return nil, fmt.Errorf("failed to download feature %s: %w", featureID, err)
		}

		// Parse the feature definition
		def, err := r.ParseFeatureDefinition(tarballPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse feature %s: %w", featureID, err)
		}

		result[featureID] = &ResolvedFeature{
			Definition: def,
			TarballPath: tarballPath,
			InstallPath: "",
		}
	}

	// Resolve dependencies using topological sort
	ordered, err := r.resolveDependencies(result)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve feature dependencies: %w", err)
	}

	// Reorder result based on dependencies
	orderedResult := make(map[string]*ResolvedFeature)
	for _, id := range ordered {
		orderedResult[id] = result[id]
	}

	return orderedResult, nil
}

// DownloadFeature downloads a feature from OCI registry and returns the OCI reference
// For now, we return the OCI reference and let the Dockerfile handle it via multi-stage build
func (r *Resolver) DownloadFeature(ctx context.Context, ref string) (string, error) {
	refParsed, err := name.ParseReference(ref)
	if err != nil {
		return "", err
	}

	// Verify the image exists by trying to fetch it
	_, err = remote.Image(refParsed, r.registryOpts...)
	if err != nil {
		return "", fmt.Errorf("failed to pull feature image %s: %w", ref, err)
	}

	// Return the full OCI reference (will be used in Dockerfile)
	return refParsed.Name(), nil
}

// ParseFeatureDefinition parses devcontainer-feature.json from a downloaded tarball
func (r *Resolver) ParseFeatureDefinition(tarballPath string) (*config.FeatureDefinition, error) {
	// For now, we'll extract and parse the label from the image
	// A more complete implementation would extract the tarball
	// This is a simplified version that shows the concept
	refParsed, err := name.ParseReference(tarballPath)
	if err != nil {
		return nil, err
	}

	img, err := remote.Image(refParsed, r.registryOpts...)
	if err != nil {
		return nil, err
	}

	cfg, err := img.ConfigFile()
	if err != nil {
		return nil, err
	}

	label := cfg.Config.Labels["devcontainer-feature"]
	if label == "" {
		// If no label, create a basic definition
		return &config.FeatureDefinition{
			ID:        "unknown",
			Name:      "Unknown Feature",
			Version:   "latest",
			DependsOn: []string{},
			Options:   make(map[string]config.OptionSpec),
		}, nil
	}

	var def config.FeatureDefinition
	if err := json.Unmarshal([]byte(label), &def); err != nil {
		return nil, fmt.Errorf("failed to parse feature definition: %w", err)
	}

	return &def, nil
}

// resolveDependencies resolves feature dependencies using topological sort
func (r *Resolver) resolveDependencies(features map[string]*ResolvedFeature) ([]string, error) {
	// Convert to config.FeatureDefinition map for topological sort
	defs := make(map[string]*config.FeatureDefinition)
	for id, f := range features {
		defs[id] = f.Definition
	}

	return config.TopologicalSort(defs)
}
