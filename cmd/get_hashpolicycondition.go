package cmd

import (
    "encoding/json"
    "fmt"
    "os"
    "text/tabwriter"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
    "dtctl/pkg/dependencytrack"
)

var ghPolicyUUID string
var ghProjectTag string

var getHashPolicyConditionCmd = &cobra.Command{
    Use:   "hashpolicycondition",
    Short: "Get hash policy condition details",
    RunE:  getHashPolicyCondition,
}

func init() {
    getHashPolicyConditionCmd.Flags().StringVar(&ghPolicyUUID, "policy-uuid", "", "UUID of the policy (optional)")
    getHashPolicyConditionCmd.Flags().StringVar(&ghProjectTag, "project-tag", "", "Filter by project tag (optional)")

    getCmd.AddCommand(getHashPolicyConditionCmd)
}

func getHashPolicyCondition(cmd *cobra.Command, args []string) error {
    cfg, err := config.GetConfig()
    if err != nil {
        return err
    }
    if cfg.CurrentContext == "" {
        return fmt.Errorf("no current context is set; use 'dtctl config use-context'")
    }

    ctx, err := config.GetContext(cfg.CurrentContext)
    if err != nil {
        return err
    }

    client := dependencytrack.NewClient(ctx.URL, ctx.Token)

    // If project-tag is given, we filter policies by projects that have this tag
    // If policy-uuid is given, we show that specific policy
    // If neither given, no action or maybe show an error
    if ghPolicyUUID == "" && ghProjectTag == "" {
        return fmt.Errorf("either --policy-uuid or --project-tag must be provided")
    }

    //var taggedProjects []dependencytrack.Project
    var taggedProjectUUIDs = make(map[string]bool)
    if ghProjectTag != "" {
        // Get projects by tag
        p, err := client.GetProjectsByTag(ghProjectTag)
        if err != nil {
            return fmt.Errorf("failed to get projects by tag: %v", err)
        }
        //taggedProjects = p
        for _, proj := range p {
            taggedProjectUUIDs[proj.UUID] = true
        }
    }

    var policies []map[string]interface{}

    if ghPolicyUUID != "" {
        // Get a single policy by UUID
        pol, err := client.GetPolicyByUUID(ghPolicyUUID)
        if err != nil {
            return fmt.Errorf("failed to get policy: %v", err)
        }
        policies = append(policies, pol)
    } else {
        // Get all policies
        allPolicies, err := client.GetPolicies()
        if err != nil {
            return fmt.Errorf("failed to get all policies: %v", err)
        }

        // Convert to []map[string]interface{} from []Policy only if needed,
        // but we previously had GetPolicies return []Policy. If it returns []Policy,
        // you must adjust code. If previously it returned []map[string]interface{}, skip this.
        // Let's assume GetPolicies() returns []Policy, we must cast to interface:
        // Actually from previous code, GetPolicies returns []Policy (struct). We need conditions.
        // Let's re-check that scenario:

        // If GetPolicies() returns []Policy (struct):
        // We need conditions and projects from it. It's easier if we had a map. Let's assume we modify GetPolicies()
        // to return []map[string]interface{} as well. If not, we must rely on a different approach.
        // For simplicity, let's assume it's returning []Policy minimal and we can't access conditions from the struct.
        // We must get each policy by UUID or store more data. To keep it simple:
        // We'll do a second step: after getting all policies (struct), we do another GetPolicyByUUID for each to get conditions.

        // We'll do that for each policy returned:
        // If your Policy struct doesn't hold enough info, you must individually query again by UUID.
        // Let's assume your Policy struct has at least Name and UUID.

        for _, polStruct := range allPolicies {
            // polStruct is type Policy {Name, UUID, Projects ...}
            // Convert to map for consistency or directly handle logic:
            // Instead of a map, let's just do another API call:
            polMap, err := client.GetPolicyByUUID(polStruct.UUID)
            if err != nil {
                return fmt.Errorf("failed to get policy by UUID %s: %v", polStruct.UUID, err)
            }
            policies = append(policies, polMap)
        }
    }

    var results [][]string

    for _, policy := range policies {
        policyName, _ := policy["name"].(string)
        policyConditions, ok := policy["policyConditions"].([]interface{})
        if !ok || len(policyConditions) == 0 {
            continue // no conditions
        }

        condition := policyConditions[0].(map[string]interface{})
        operator, _ := condition["operator"].(string)
        valueStr, _ := condition["value"].(string)

        var valObj map[string]string
        if err := json.Unmarshal([]byte(valueStr), &valObj); err != nil {
            return fmt.Errorf("failed to parse value: %v", err)
        }

        algorithm := valObj["algorithm"]
        algorithmValue := valObj["value"]

        projects, ok := policy["projects"].([]interface{})
        if !ok {
            // No projects
            // If we have a project-tag filter, then no match since no projects
            if ghProjectTag == "" {
                // Print condition anyway
                results = append(results, []string{policyName, "", operator, algorithm, algorithmValue})
            }
            continue
        }

        for _, p := range projects {
            proj := p.(map[string]interface{})
            projectName, _ := proj["name"].(string)
            projectUUID, _ := proj["uuid"].(string)

            if ghProjectTag != "" {
                // Filter only if this project is in taggedProjectUUIDs
                if !taggedProjectUUIDs[projectUUID] {
                    continue
                }
            }

            results = append(results, []string{policyName, projectName, operator, algorithm, algorithmValue})
        }
    }

    if len(results) == 0 {
        fmt.Println("No hash policy conditions found.")
        return nil
    }

    printHashPolicyCondition(results)
    return nil
}

func printHashPolicyCondition(results [][]string) {
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "Policy Name\tProject Name\tOperator\tAlgorithm\tAlgorithm Value")
    fmt.Fprintln(w, "-----------\t------------\t--------\t---------\t--------------")
    for _, row := range results {
        fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", row[0], row[1], row[2], row[3], row[4])
    }
    w.Flush()
}
