package cmd

import (
    "fmt"
    "os"
    "text/tabwriter"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
    "dtctl/pkg/dependencytrack"
)

var showProjects bool

func init() {
    getPoliciesCmd.Flags().BoolVar(&showProjects, "show-projects", false, "Show associated projects for each policy")
}

var getPoliciesCmd = &cobra.Command{
    Use:   "policies",
    Short: "Get policies",
    RunE:  getPolicies,
}

func getPolicies(cmd *cobra.Command, args []string) error {
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
    client := dependencytrack.NewClient(ctx.URL, ctx.Token)

    policies, err := client.GetPolicies()
    if err != nil {
        return err
    }

    if len(policies) == 0 {
        fmt.Println("No policies found.")
        return nil
    }

    // Display data in table format
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    if showProjects {
        fmt.Fprintln(w, "POLICY NAME\tPOLICY UUID\tPROJECTS")
        fmt.Fprintln(w, "-----------\t-----------\t--------")
        for _, policy := range policies {
            var projectNames []string
            for _, project := range policy.Projects {
                projectNames = append(projectNames, project.Name)
            }
            fmt.Fprintf(w, "%s\t%s\t%s\n", policy.Name, policy.UUID, joinStrings(projectNames, ", "))
        }
    } else {
        fmt.Fprintln(w, "POLICY NAME\tPOLICY UUID")
        fmt.Fprintln(w, "-----------\t-----------")
        for _, policy := range policies {
            fmt.Fprintf(w, "%s\t%s\n", policy.Name, policy.UUID)
        }
    }
    w.Flush()

    return nil
}

// Helper function to join strings
func joinStrings(strs []string, sep string) string {
    result := ""
    for i, s := range strs {
        result += s
        if i < len(strs)-1 {
            result += sep
        }
    }
    return result
}
