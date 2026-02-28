package feature

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
func (r *Resolver) ResolveLocalFeatures(dir string) (map[string]*config.FeatureDefinition, error) {
	featuresDir := filepath.Join(dir, ".devcontainer", "features")

	entries, err := os.ReadDir(featuresDir)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]*config.FeatureDefinition), nil
		}
		return nil, err
	}

	result := make(map[string]*config.FeatureDefinition)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		featurePath := filepath.Join(featuresDir, entry.Name())
		def, err := config.ParseFeatureDefinition(featurePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse feature %s: %w", entry.Name(), err)
		}

		result[def.ID] = def
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
