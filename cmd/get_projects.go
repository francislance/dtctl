package cmd

import (
    "fmt"
    "os"
    "text/tabwriter"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
    "dtctl/pkg/dependencytrack"
)

var tag string

func init() {
    getCmd.AddCommand(getProjectsCmd)
    getProjectsCmd.Flags().StringVar(&tag, "tag", "", "Filter projects by tag")
}

var getProjectsCmd = &cobra.Command{
    Use:   "projects",
    Short: "Get all projects",
    RunE: func(cmd *cobra.Command, args []string) error {
        cfg, err := config.GetConfig()
        if err != nil {
            return err
        }
        if cfg.CurrentContext == "" {
            return fmt.Errorf("no current context is set; use 'dtctl config use-context' to set one")
        }

        // Change this line:
        // ctx, err := cfg.GetCurrentContext()
        // To:
        ctx, err := config.GetCurrentContext()
        if err != nil {
            return err
        }

        client := dependencytrack.NewClient(ctx.URL, ctx.Token)

        var projects []dependencytrack.Project

        if tag != "" {
            projects, err = client.GetProjectsByTag(tag)
            if err != nil {
                return err
            }
        } else {
            projects, err = client.GetProjects()
            if err != nil {
                return err
            }
        }

        if len(projects) == 0 {
            fmt.Println("No projects found.")
            return nil
        }

        // Use tabwriter for formatted output
        w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
        // Print header
        fmt.Fprintln(w, "NAME\tUUID")
        // Print separator line
        fmt.Fprintln(w, "----\t----")
        // Print project data
        for _, project := range projects {
            fmt.Fprintf(w, "%s\t%s\n", project.Name, project.UUID)
        }
        // Flush the writer to ensure all data is written
        w.Flush()

        return nil
    },
}
