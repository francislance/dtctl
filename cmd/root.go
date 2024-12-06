package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:     "dtctl",
    Short:   "dtctl is a CLI tool for interacting with Dependency-Track",
    Version: "", // Will set the version in init()
}

// Execute executes the root command.
func Execute() error {
    return rootCmd.Execute()
}

func init() {
    // Set the version using the GetVersion function from version.go
    rootCmd.Version = GetVersion()

    // Customize the version output format
    rootCmd.SetVersionTemplate(fmt.Sprintf("dtctl %s\n", GetVersion()))

    // Add subcommands
    rootCmd.AddCommand(configCmd)
    rootCmd.AddCommand(getCmd)
}
