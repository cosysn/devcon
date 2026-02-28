package main

import (
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up <dir>",
	Short: "Start devcontainer for local testing",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
