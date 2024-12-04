package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
)

var (
    url   string
    token string
)

func init() {
    configCmd.AddCommand(addContextCmd)
    addContextCmd.Flags().StringVar(&url, "url", "", "Dependency-Track server URL")
    addContextCmd.Flags().StringVar(&token, "token", "", "API token")
    addContextCmd.MarkFlagRequired("url")
    addContextCmd.MarkFlagRequired("token")
}

var addContextCmd = &cobra.Command{
    Use:   "add-context NAME",
    Short: "Add a new context",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        name := args[0]
        ctx := config.Context{
            Name:  name,
            URL:   url,
            Token: token,
        }
        if err := config.AddContext(ctx); err != nil {
            return err
        }
        fmt.Printf("Context '%s' added successfully.\n", name)
        return nil
    },
}
