package cmd

import (
    "github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
    Use:   "get",
    Short: "Retrieve resources",
}

func init() {
    rootCmd.AddCommand(getCmd)
    getCmd.AddCommand(getProjectsCmd)
    getCmd.AddCommand(getComponentsCmd)
    getCmd.AddCommand(getPoliciesCmd)
}
