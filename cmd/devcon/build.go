package main

import (
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build <dir>",
	Short: "Build devcontainer image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	buildCmd.Flags().String("provider", "docker", "Build provider (docker)")
	rootCmd.AddCommand(buildCmd)
}
