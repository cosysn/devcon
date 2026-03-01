package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devcon/cli/internal/builder"
	"github.com/devcon/cli/pkg/config"
	"github.com/devcon/cli/pkg/feature"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build <dir>",
	Short: "Build devcontainer image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		provider, _ := cmd.Flags().GetString("provider")

		// Validate directory
		stat, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("failed to access directory: %w", err)
		}
		if !stat.IsDir() {
			return fmt.Errorf("path is not a directory: %s", dir)
		}

		// Parse config
		cfg, err := config.ParseDevcontainer(dir)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		// Resolve extends
		cfg, err = config.ResolveExtends(dir, cfg)
		if err != nil {
			return fmt.Errorf("failed to resolve extends: %w", err)
		}

		// Validate we have either image or dockerfile
		if cfg.Image == "" && cfg.Dockerfile == "" {
			return fmt.Errorf("either image or dockerFile must be specified in devcontainer.json")
		}

		// Resolve and download features if any are specified
		var resolvedFeatures map[string]*feature.ResolvedFeature
		if len(cfg.Features) > 0 {
			fmt.Println("Resolving features...")

			resolver := feature.NewResolver()
			resolvedFeatures, err = resolver.ResolveAndDownload(context.Background(), dir, cfg.Features)
			if err != nil {
				return fmt.Errorf("failed to resolve features: %w", err)
			}

			fmt.Printf("Resolved %d features\n", len(resolvedFeatures))

			// Generate Dockerfile with features if we have features and an image
			if cfg.Image != "" && len(resolvedFeatures) > 0 {
				dockerfileContent, err := feature.GenerateDockerfile(cfg.Image, resolvedFeatures)
				if err != nil {
					return fmt.Errorf("failed to generate Dockerfile: %w", err)
				}

				// Write generated Dockerfile to .devcontainer directory
				generatedDockerfile := filepath.Join(dir, ".devcontainer", "generated.Dockerfile")
				if err := os.WriteFile(generatedDockerfile, []byte(dockerfileContent), 0644); err != nil {
					return fmt.Errorf("failed to write generated Dockerfile: %w", err)
				}

				// Create a combined Dockerfile that uses the base and includes features
				combinedDockerfile := filepath.Join(dir, "Dockerfile.with-features")
				combinedContent := fmt.Sprintf("FROM %s\n\n", cfg.Image) + dockerfileContent
				if err := os.WriteFile(combinedDockerfile, []byte(combinedContent), 0644); err != nil {
					return fmt.Errorf("failed to write combined Dockerfile: %w", err)
				}

				// Use the combined Dockerfile for build
				cfg.Dockerfile = "Dockerfile.with-features"
				cfg.Image = "" // Use Dockerfile, not image directly
			}
		}

		// Create builder
		b, err := builder.NewBuilder(provider)
		if err != nil {
			return fmt.Errorf("failed to create builder: %w", err)
		}

		// Build
		spec := builder.Spec{
			ContextDir: dir,
			Image:      cfg.Image,
			Dockerfile: cfg.Dockerfile,
			Features:   cfg.Features,
			Env:        cfg.Env,
		}

		imageID, err := b.Build(context.Background(), spec)
		if err != nil {
			return fmt.Errorf("build failed: %w", err)
		}

		fmt.Println("Image built:", imageID)
		return nil
	},
}

func init() {
	buildCmd.Flags().String("provider", "docker", "Build provider (docker)")
	rootCmd.AddCommand(buildCmd)
}
