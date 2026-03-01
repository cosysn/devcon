package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devcon/cli/internal/builder"
	"github.com/devcon/cli/pkg/config"
	"github.com/devcon/cli/pkg/errors"
	"github.com/devcon/cli/pkg/feature"
	"github.com/devcon/cli/pkg/output"
	"github.com/devcon/cli/pkg/progress"
	"github.com/spf13/cobra"
)

// buildLoggerAdapter adapts output to builder.BuildLogger
type upBuildLoggerAdapter struct {
	output output.Output
}

func (l *upBuildLoggerAdapter) Write(line string) {
	if line != "" {
		l.output.Verbose(line)
	}
}

var upCmd = &cobra.Command{
	Use:   "up <dir>",
	Short: "Start devcontainer for local testing",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		provider, _ := cmd.Flags().GetString("provider")
		out := output.GetGlobalOutput()

		// Validate directory
		stat, err := os.Stat(dir)
		if err != nil {
			return errors.NewEnhancedError(errors.ErrCodeUnknown,
				fmt.Sprintf("failed to access directory: %v", err), err, nil)
		}
		if !stat.IsDir() {
			return errors.NewEnhancedError(errors.ErrCodeUnknown,
				fmt.Sprintf("path is not a directory: %s", dir), nil, nil)
		}

		out.Verbosef("Working directory: %s", dir)

		// Parse config
		configPath := filepath.Join(dir, ".devcontainer", "devcontainer.json")
		out.Verbosef("Reading config: %s", configPath)

		cfg, err := config.ParseDevcontainer(dir)
		if err != nil {
			return errors.NewEnhancedError(errors.ErrCodeConfigParse,
				fmt.Sprintf("failed to parse config: %v", err), err,
				&errors.Suggestion{
					Text:   "Check your devcontainer.json for syntax errors",
					Action: "Use a JSON validator or IDE with JSON schema support",
				})
		}

		out.Verbose("Config parsed successfully")
		out.Verbosef("  Image: %s", cfg.Image)
		out.Verbosef("  Dockerfile: %s", cfg.Dockerfile)
		if len(cfg.Features) > 0 {
			out.Verbosef("  Features: %v", cfg.Features)
		}
		if len(cfg.ContainerEnv) > 0 {
			out.Verbosef("  Environment: %v", cfg.ContainerEnv)
		}

		// Resolve extends
		if cfg.Extends != "" {
			out.Verbosef("Resolving extends: %s", cfg.Extends)
		}
		cfg, err = config.ResolveExtends(dir, cfg)
		if err != nil {
			return errors.NewEnhancedError(errors.ErrCodeConfigExtend,
				fmt.Sprintf("failed to resolve extends: %v", err), err,
				&errors.Suggestion{
					Text:   "Check your extends configuration",
					Action: "Ensure the extended config file exists and is valid",
				})
		}
		if cfg.Extends != "" {
			out.Verbose("Extends resolved successfully")
		}

		// Validate we have either image or dockerfile
		if cfg.Image == "" && cfg.Dockerfile == "" {
			return errors.NewMissingImageOrDockerfileError()
		}

		// Resolve and download features if any are specified
		if len(cfg.Features) > 0 {
			p := progress.NewSpinner("Resolving features")
			p.Start()

			out.Verbosef("Resolving %d features...", len(cfg.Features))
			for featID := range cfg.Features {
				out.Verbosef("  - Resolving: %s", featID)
			}

			resolver := feature.NewResolver()
			resolvedFeatures, err := resolver.ResolveAndDownload(context.Background(), dir, cfg.Features)
			if err != nil {
				p.Stop("failed")
				return errors.NewEnhancedError(errors.ErrCodeFeatureResolve,
					fmt.Sprintf("failed to resolve features: %v", err), err,
					&errors.Suggestion{
						Text:   "Check the feature IDs and versions",
						Action: "Ensure features exist in the specified registry",
					})
			}

			p.Stop(fmt.Sprintf("Resolved %d features", len(resolvedFeatures)))
			out.Verbosef("Features resolved:")
			for name, f := range resolvedFeatures {
				if f.Definition != nil {
					out.Verbosef("  - %s: version=%s, id=%s", name, f.Definition.Version, f.Definition.ID)
				}
				if f.TarballPath != "" {
					out.Verbosef("    Tarball: %s", f.TarballPath)
				}
				if f.InstallPath != "" {
					out.Verbosef("    Install path: %s", f.InstallPath)
				}
			}

			// Generate Dockerfile with features if we have features and an image
			if cfg.Image != "" && len(resolvedFeatures) > 0 {
				out.Verbose("Generating Dockerfile with features...")
				dockerfileContent, err := feature.GenerateDockerfile(cfg.Image, resolvedFeatures)
				if err != nil {
					return fmt.Errorf("failed to generate Dockerfile: %w", err)
				}

				// Write generated Dockerfile to .devcontainer directory
				generatedDockerfile := filepath.Join(dir, ".devcontainer", "generated.Dockerfile")
				if err := os.WriteFile(generatedDockerfile, []byte(dockerfileContent), 0644); err != nil {
					return fmt.Errorf("failed to write generated Dockerfile: %w", err)
				}
				out.Verbosef("Written: %s", generatedDockerfile)

				// Create a combined Dockerfile that uses the base and includes features
				combinedDockerfile := filepath.Join(dir, "Dockerfile.with-features")
				combinedContent := fmt.Sprintf("FROM %s\n\n", cfg.Image) + dockerfileContent
				if err := os.WriteFile(combinedDockerfile, []byte(combinedContent), 0644); err != nil {
					return fmt.Errorf("failed to write combined Dockerfile: %w", err)
				}
				out.Verbosef("Written: %s", combinedDockerfile)

				// Use the combined Dockerfile for build
				cfg.Dockerfile = "Dockerfile.with-features"
				cfg.Image = "" // Use Dockerfile, not image directly
			}
		}

		// Create builder with logger
		out.Verbosef("Initializing builder: %s", provider)
		logger := &upBuildLoggerAdapter{output: out}
		b, err := builder.NewBuilderWithLogger(provider, logger)
		if err != nil {
			return errors.NewEnhancedError(errors.ErrCodeBuilderInit,
				fmt.Sprintf("failed to create builder: %v", err), err,
				&errors.Suggestion{
					Text:   "Check the build provider configuration",
					Action: "Ensure the provider is installed and accessible",
				})
		}

		// Build first
		// Convert lifecycle commands to strings for backward compatibility
		onCreateCmd := cfg.OnCreateCommand.ToString()
		postCreateCmd := cfg.PostCreateCommand.ToString()
		postStartCmd := cfg.PostStartCommand.ToString()

		spec := builder.Spec{
			ContextDir:         dir,
			Image:             cfg.Image,
			Dockerfile:         cfg.Dockerfile,
			Features:          cfg.Features,
			Env:               cfg.ContainerEnv,
			Mounts:            cfg.Mounts,
			Privileged:        cfg.Privileged != nil && *cfg.Privileged,
			OnCreateCommand:   onCreateCmd,
			PostCreateCommand: postCreateCmd,
			PostStartCommand:  postStartCmd,
		}

		out.Verbosef("Build context: %s", spec.ContextDir)
		if spec.Dockerfile != "" {
			out.Verbosef("Using Dockerfile: %s", spec.Dockerfile)
		}
		if spec.Image != "" {
			out.Verbosef("Using image: %s", spec.Image)
		}
		if spec.OnCreateCommand != "" {
			out.Verbosef("onCreateCommand: %s", spec.OnCreateCommand)
		}
		if spec.PostCreateCommand != "" {
			out.Verbosef("postCreateCommand: %s", spec.PostCreateCommand)
		}
		if spec.PostStartCommand != "" {
			out.Verbosef("postStartCommand: %s", spec.PostStartCommand)
		}

		p := progress.NewSpinner("Building image")
		p.Start()
		out.Verbose("Starting Docker build...")

		imageID, err := b.Build(context.Background(), spec)
		if err != nil {
			p.Stop("failed")
			return errors.NewBuildFailedError(err)
		}

		p.Stop("Image built")
		out.Verbosef("Build completed. Image ID: %s", imageID)

		// Start container
		spec.Image = imageID
		out.Verbosef("Starting container from image: %s", imageID)
		out.StartProgress("Starting container")
		if err := b.Up(context.Background(), spec); err != nil {
			out.StopProgress("failed")
			return errors.NewEnhancedError(errors.ErrCodeContainerStart,
				fmt.Sprintf("failed to start container: %v", err), err,
				&errors.Suggestion{
					Text:   "Check container logs for more details",
					Action: "Run 'docker ps' to check container status",
				})
		}
		out.StopProgress("Container started")

		out.Success("Devcontainer is running!")
		return nil
	},
}

func init() {
	upCmd.Flags().String("provider", "docker", "Build provider (docker)")
	rootCmd.AddCommand(upCmd)
}
