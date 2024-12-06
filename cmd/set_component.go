package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
    "dtctl/pkg/dependencytrack"
)

var (
    componentUUID string
    newSHA256     string
)

// setComponentCmd represents the set component command
var setComponentCmd = &cobra.Command{
    Use:   "component",
    Short: "Set or update a component's fields",
    RunE:  setComponent,
}

func init() {
    setComponentCmd.Flags().StringVarP(&componentUUID, "uuid", "u", "", "UUID of the component (required)")
    setComponentCmd.Flags().StringVar(&newSHA256, "field-sha256", "", "New SHA256 value for the component (required)")
    setComponentCmd.MarkFlagRequired("uuid")
    setComponentCmd.MarkFlagRequired("field-sha256")
    setCmd.AddCommand(setComponentCmd)
}

// setComponent handles the execution of the set component command
func setComponent(cmd *cobra.Command, args []string) error {
    // Retrieve configuration
    cfg, err := config.GetConfig()
    if err != nil {
        return err
    }
    if cfg.CurrentContext == "" {
        return fmt.Errorf("no current context is set; use 'dtctl config use-context' to set one")
    }
    ctx, err := config.GetCurrentContext()
    if err != nil {
        return err
    }

    // Initialize the Dependency-Track client
    client := dependencytrack.NewClient(ctx.URL, ctx.Token)

    // Update the sha256 field
    err = client.UpdateComponentSHA256(componentUUID, newSHA256)
    if err != nil {
        return fmt.Errorf("failed to update component: %v", err)
    }

    fmt.Println("Component sha256 updated successfully.")
    return nil
}
