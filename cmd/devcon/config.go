package main

import (
	"fmt"
	"os"

	"github.com/devcon/cli/pkg/config"
	"github.com/devcon/cli/pkg/output"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config <dir>",
	Short: "Validate and show devcontainer config",
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

		// Resolve extends
		cfg, err = config.ResolveExtends(dir, cfg)
		if err != nil {
			return fmt.Errorf("failed to resolve extends: %w", err)
		}

		// Output - use global output mode
		// Note: The global --output flag takes precedence
		out.Success("Config validated successfully")
		out.Printf("Image: %s\n", cfg.Image)
		out.Printf("Dockerfile: %s\n", cfg.Dockerfile)
		out.Printf("Features: %v\n", cfg.Features)
		out.Printf("Env: %v\n", cfg.Env)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
