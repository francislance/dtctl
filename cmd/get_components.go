package cmd

import (
    "fmt"
    "os"
    "text/tabwriter"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
    "dtctl/pkg/dependencytrack"
)

var componentTag string

func init() {
    getCmd.AddCommand(getComponentsCmd)
    getComponentsCmd.Flags().StringVar(&componentTag, "tag", "", "Filter components by project tag (required)")
    getComponentsCmd.MarkFlagRequired("tag")
}

var getComponentsCmd = &cobra.Command{
    Use:   "components",
    Short: "Get components based on project tag",
    RunE:  getComponents,
}

func getComponents(cmd *cobra.Command, args []string) error {
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

    projects, err := client.GetProjectsByTag(componentTag)
    if err != nil {
        return err
    }

    if len(projects) == 0 {
        fmt.Println("No projects found with the specified tag.")
        return nil
    }

    // Prepare data for display
    type ComponentInfo struct {
        ProjectName    string
        ProjectUUID    string
        ComponentName  string
        ComponentUUID  string
        Sha256         string
    }
    var components []ComponentInfo

    for _, project := range projects {
        projectComponents, err := client.GetComponentsByProjectUUID(project.UUID)
        if err != nil {
            return err
        }
        for _, component := range projectComponents {
            components = append(components, ComponentInfo{
                ProjectName:   project.Name,
                ProjectUUID:   project.UUID,
                ComponentName: component.Name,
                ComponentUUID: component.UUID,
                Sha256:        component.Sha256,
            })
        }
    }

    if len(components) == 0 {
        fmt.Println("No components found.")
        return nil
    }

    // Display data in table format
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "PROJECT NAME\tPROJECT UUID\tCOMPONENT NAME\tCOMPONENT UUID\tSHA256")
    fmt.Fprintln(w, "------------\t------------\t--------------\t--------------\t------")
    for _, comp := range components {
        fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", comp.ProjectName, comp.ProjectUUID, comp.ComponentName, comp.ComponentUUID, comp.Sha256)
    }
    w.Flush()

    return nil
}
