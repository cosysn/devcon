package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/devcon/cli/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config <dir>",
	Short: "Validate and show devcontainer config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		output, _ := cmd.Flags().GetString("output")

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

		// Output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(cfg); err != nil {
				return fmt.Errorf("failed to encode config: %w", err)
			}
		} else {
			fmt.Printf("Image: %s\n", cfg.Image)
			fmt.Printf("Dockerfile: %s\n", cfg.Dockerfile)
			fmt.Printf("Features: %v\n", cfg.Features)
			fmt.Printf("Env: %v\n", cfg.Env)
		}

		return nil
	},
}

func init() {
	configCmd.Flags().String("output", "text", "Output format (text, json)")
	rootCmd.AddCommand(configCmd)
}
