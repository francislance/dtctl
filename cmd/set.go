package cmd

import (
    "github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
    Use:   "set",
    Short: "Set or update resources",
}

func init() {
    rootCmd.AddCommand(setCmd)
    setCmd.AddCommand(setComponentCmd)
}
