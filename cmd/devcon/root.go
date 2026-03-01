package main

import (
	"os"

	"github.com/devcon/cli/pkg/output"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devcon",
	Short: "devcontainer CLI tool",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		output.ApplyOutputSettings(cmd, cmd.Name())
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringP("output", "o", "text", "Output format (text, json)")
	rootCmd.PersistentFlags().Bool("quiet", false, "Suppress non-essential output")

	rootCmd.AddCommand(featuresCmd)
}
