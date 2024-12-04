package cmd

import (
    "github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Manage dtctl configurations",
}

func init() {
    // Subcommands will be added in their respective files
}
