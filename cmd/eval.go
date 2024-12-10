// cmd/eval.go
package cmd

import (
    "github.com/spf13/cobra"
)

var evalCmd = &cobra.Command{
    Use:   "eval",
    Short: "Evaluate conditions or policies",
}

func init() {
    rootCmd.AddCommand(evalCmd)
}
