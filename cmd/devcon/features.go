package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devcon/cli/pkg/feature"
	"github.com/spf13/cobra"
)

var featuresCmd = &cobra.Command{
	Use:   "features",
	Short: "Manage devcontainer features",
}

var featuresPackageCmd = &cobra.Command{
	Use:   "package <dir>",
	Short: "Package a Feature as OCI artifact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		output, _ := cmd.Flags().GetString("output")

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

		if output == "" {
			output = filepath.Join(dir, "feature.tar.gz")
		}

		if err := feature.PackageFeature(dir, output); err != nil {
			return fmt.Errorf("failed to package feature: %w", err)
		}

		fmt.Println("Feature packaged:", output)
		return nil
	},
}

var featuresPublishCmd = &cobra.Command{
	Use:   "publish <dir>",
	Short: "Publish a Feature to OCI registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	featuresPackageCmd.Flags().String("output", "", "Output path")
	featuresCmd.AddCommand(featuresPackageCmd)
	featuresCmd.AddCommand(featuresPublishCmd)
}
