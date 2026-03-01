package main

import (
	"fmt"
	"os"

	"github.com/devcon/cli/pkg/config"
	"github.com/devcon/cli/pkg/feature"
	"github.com/devcon/cli/pkg/output"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <dir>",
	Short: "Inspect parsed config and feature dependencies",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		out := output.GetGlobalOutput()

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

		cfg, err = config.ResolveExtends(dir, cfg)
		if err != nil {
			return fmt.Errorf("failed to resolve extends: %w", err)
		}

		// Resolve features
		resolver := feature.NewResolver()
		localFeatures, err := resolver.ResolveLocalFeatures(dir)
		if err != nil {
			return fmt.Errorf("failed to resolve local features: %w", err)
		}

		// Output
		out.Println("=== Config ===")
		out.Printf("Image: %s\n", cfg.Image)
		out.Printf("Dockerfile: %s\n", cfg.Dockerfile)
		out.Printf("Features: %v\n", cfg.Features)

		out.Println("\n=== Local Features ===")
		if len(localFeatures) == 0 {
			out.Println("No local features found")
		} else {
			for id, f := range localFeatures {
				out.Printf("- %s (version: %s)\n", id, f.Version)
				if len(f.DependsOn) > 0 {
					out.Printf("  dependsOn: %v\n", f.DependsOn)
				}
			}
		}

		out.Success("Inspection complete")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
