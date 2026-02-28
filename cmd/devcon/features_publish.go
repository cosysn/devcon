package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devcon/cli/pkg/feature"
	"github.com/spf13/cobra"
)

var featuresPublishCmd = &cobra.Command{
	Use:   "publish <dir>",
	Short: "Publish a Feature to OCI registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		reg, _ := cmd.Flags().GetString("reg")

		if reg == "" {
			return fmt.Errorf("--reg is required")
		}

		// Validate input directory
		stat, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("failed to access directory: %w", err)
		}
		if !stat.IsDir() {
			return fmt.Errorf("path is not a directory: %s", dir)
		}

		// Check for required files
		featureJSON := filepath.Join(dir, "devcontainer-feature.json")
		if _, err := os.Stat(featureJSON); os.IsNotExist(err) {
			return fmt.Errorf("devcontainer-feature.json not found in %s", dir)
		}

		ctx := context.Background()
		if err := feature.PublishFeature(ctx, dir, reg); err != nil {
			return fmt.Errorf("failed to publish feature: %w", err)
		}

		fmt.Println("Feature published:", reg)
		return nil
	},
}

func init() {
	featuresPublishCmd.Flags().String("reg", "", "Registry URL (e.g., ghcr.io/user/feature)")
	featuresCmd.AddCommand(featuresPublishCmd)
}
