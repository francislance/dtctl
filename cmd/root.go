package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "dtctl",
    Short: "dtctl is a CLI tool for interacting with Dependency-Track",
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    // Add subcommands
    rootCmd.AddCommand(configCmd)
    rootCmd.AddCommand(getCmd)
}
