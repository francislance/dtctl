package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
)

func init() {
    configCmd.AddCommand(getContextsCmd)
}

var getContextsCmd = &cobra.Command{
    Use:   "get-contexts",
    Short: "List all contexts",
    RunE: func(cmd *cobra.Command, args []string) error {
        cfg, err := config.GetConfig()
        if err != nil {
            return err
        }
        fmt.Println("Available contexts:")
        for _, ctx := range cfg.Contexts {
            current := " "
            if ctx.Name == cfg.CurrentContext {
                current = "*"
            }
            fmt.Printf("%s %s\n", current, ctx.Name)
        }
        return nil
    },
}
