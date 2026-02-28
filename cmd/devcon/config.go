package main

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config <dir>",
	Short: "Validate and show devcontainer config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
