package main

import (
	"context"
	"fmt"
	"os"

	"github.com/devcon/cli/internal/builder"
	"github.com/devcon/cli/pkg/config"
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

		// Create builder
		b, err := builder.NewBuilder(provider)
		if err != nil {
			return fmt.Errorf("failed to create builder: %w", err)
		}

		// Build
		spec := builder.Spec{
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
