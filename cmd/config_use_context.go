package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
)

func init() {
    configCmd.AddCommand(useContextCmd)
}

var useContextCmd = &cobra.Command{
    Use:   "use-context NAME",
    Short: "Set the current context",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        name := args[0]
        if err := config.UseContext(name); err != nil {
            return err
        }
        fmt.Printf("Switched to context '%s'.\n", name)
        return nil
    },
}
