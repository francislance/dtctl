package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
)

var (
    editURL   string
    editToken string
)

func init() {
    configCmd.AddCommand(editContextCmd)
    editContextCmd.Flags().StringVar(&editURL, "url", "", "New Dependency-Track server URL")
    editContextCmd.Flags().StringVar(&editToken, "token", "", "New API token")
}

var editContextCmd = &cobra.Command{
    Use:   "edit-context NAME",
    Short: "Edit an existing context",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        name := args[0]
        if editURL == "" && editToken == "" {
            return fmt.Errorf("no changes specified; use --url and/or --token to modify the context")
        }

        ctx, err := config.GetContext(name)
        if err != nil {
            return err
        }

        if editURL != "" {
            ctx.URL = editURL
        }
        if editToken != "" {
            ctx.Token = editToken
        }

        // Update the context in the configuration
        err = config.UpdateContext(*ctx)
        if err != nil {
            return err
        }

        fmt.Printf("Context '%s' updated successfully.\n", name)
        return nil
    },
}
