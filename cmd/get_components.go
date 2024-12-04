package cmd

import (
    "fmt"
    "os"
    "strings"
    "text/tabwriter"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
    "dtctl/pkg/dependencytrack"
)

var (
    componentTag string
    showFields   string
)

func init() {
    getCmd.AddCommand(getComponentsCmd)
    getComponentsCmd.Flags().StringVar(&componentTag, "tag", "", "Filter components by project tag (optional)")
    getComponentsCmd.Flags().StringVar(&showFields, "show-fields", "", "Comma-separated list of additional fields to display (available: projectname, projectuuid, sha256, sha1, md5)")
}

var getComponentsCmd = &cobra.Command{
    Use:   "components",
    Short: "Get components",
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

    var projects []dependencytrack.Project

    if componentTag != "" {
        projects, err = client.GetProjectsByTag(componentTag)
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

    // Prepare data for display
    type ComponentInfo struct {
        ComponentName string
        ComponentUUID string
        ProjectName   string
        ProjectUUID   string
        Sha256        string
        Sha1          string
        Md5           string
    }
    var components []ComponentInfo

    for _, project := range projects {
        projectComponents, err := client.GetComponentsByProjectUUID(project.UUID)
        if err != nil {
            return err
        }
        for _, component := range projectComponents {
            components = append(components, ComponentInfo{
                ComponentName: component.Name,
                ComponentUUID: component.UUID,
                ProjectName:   project.Name,
                ProjectUUID:   project.UUID,
                Sha256:        component.Sha256,
                Sha1:          component.Sha1,
                Md5:           component.Md5,
            })
        }
    }

    if len(components) == 0 {
        fmt.Println("No components found.")
        return nil
    }

    // Default headers and extractors
    headers := []string{"COMPONENT NAME", "COMPONENT UUID"}
    extractors := []func(ComponentInfo) string{
        func(ci ComponentInfo) string { return ci.ComponentName },
        func(ci ComponentInfo) string { return ci.ComponentUUID },
    }

    // Parse --show-fields flag
    if showFields != "" {
        fields := parseFields(showFields)
        for _, field := range fields {
            switch field {
            case "projectname":
                headers = append(headers, "PROJECT NAME")
                extractors = append(extractors, func(ci ComponentInfo) string { return ci.ProjectName })
            case "projectuuid":
                headers = append(headers, "PROJECT UUID")
                extractors = append(extractors, func(ci ComponentInfo) string { return ci.ProjectUUID })
            case "sha256":
                headers = append(headers, "SHA256")
                extractors = append(extractors, func(ci ComponentInfo) string { return ci.Sha256 })
            case "sha1":
                headers = append(headers, "SHA1")
                extractors = append(extractors, func(ci ComponentInfo) string { return ci.Sha1 })
            case "md5":
                headers = append(headers, "MD5")
                extractors = append(extractors, func(ci ComponentInfo) string { return ci.Md5 })
            default:
                return fmt.Errorf("invalid field: %s", field)
            }
        }
    }

    // Prepare the table
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

    // Print headers
    fmt.Fprintln(w, strings.Join(headers, "\t"))

    // Print separator
    separator := make([]string, len(headers))
    for i := range separator {
        separator[i] = strings.Repeat("-", len(headers[i]))
    }
    fmt.Fprintln(w, strings.Join(separator, "\t"))

    // Print rows
    for _, comp := range components {
        row := make([]string, len(extractors))
        for i, extract := range extractors {
            row[i] = extract(comp)
        }
        fmt.Fprintln(w, strings.Join(row, "\t"))
    }

    w.Flush()

    return nil
}

// Helper function to parse and normalize the show-fields input
func parseFields(input string) []string {
    fields := strings.Split(input, ",")
    for i, field := range fields {
        fields[i] = strings.TrimSpace(strings.ToLower(field))
    }
    return fields
}
