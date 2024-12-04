package cmd

import (
    "encoding/json"
    "fmt"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
    "gopkg.in/yaml.v2"
)

var outputFormat string

func init() {
    configCmd.AddCommand(getContextCmd)
    getContextCmd.Flags().StringVarP(&outputFormat, "output", "o", "yaml", "Output format (yaml or json)")
}

var getContextCmd = &cobra.Command{
    Use:   "get-context NAME",
    Short: "Display details of a context",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        name := args[0]
        ctx, err := config.GetContext(name)
        if err != nil {
            return err
        }
        switch outputFormat {
        case "yaml":
            data, err := yaml.Marshal(ctx)
            if err != nil {
                return err
            }
            fmt.Println(string(data))
        case "json":
            data, err := json.MarshalIndent(ctx, "", "  ")
            if err != nil {
                return err
            }
            fmt.Println(string(data))
        default:
            return fmt.Errorf("unsupported output format: %s", outputFormat)
        }
        return nil
    },
}
