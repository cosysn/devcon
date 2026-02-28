package main

import "github.com/spf13/cobra"

var featuresCmd = &cobra.Command{
	Use:   "features",
	Short: "Manage devcontainer features",
}

var featuresPackageCmd = &cobra.Command{
	Use:   "package <dir>",
	Short: "Package a Feature as OCI artifact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil // TODO: implement
	},
}

var featuresPublishCmd = &cobra.Command{
	Use:   "publish <dir>",
	Short: "Publish a Feature to OCI registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil // TODO: implement
	},
}

func init() {
	featuresCmd.AddCommand(featuresPackageCmd)
	featuresCmd.AddCommand(featuresPublishCmd)
}
