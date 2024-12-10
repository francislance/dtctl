package cmd

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "text/tabwriter"

    "github.com/spf13/cobra"
    "dtctl/pkg/config"
    "dtctl/pkg/dependencytrack"
)

var evalPolicyUUID string

var evalPolicyCmd = &cobra.Command{
    Use:   "policy",
    Short: "Evaluate if a policy is violated",
    RunE:  evalPolicy,
}

func init() {
    evalPolicyCmd.Flags().StringVar(&evalPolicyUUID, "uuid", "", "UUID of the policy (required)")
    evalPolicyCmd.MarkFlagRequired("uuid")
    evalCmd.AddCommand(evalPolicyCmd)
}

func evalPolicy(cmd *cobra.Command, args []string) error {
    // Retrieve config
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

    // Initialize client
    client := dependencytrack.NewClient(ctx.URL, ctx.Token)

    // Fetch the policy
    policy, err := client.GetPolicyByUUID(evalPolicyUUID)
    if err != nil {
        return fmt.Errorf("failed to get policy: %v", err)
    }

    policyName, _ := policy["name"].(string)
    policyConditions, ok := policy["policyConditions"].([]interface{})
    if !ok || len(policyConditions) == 0 {
        fmt.Println("No policy conditions found. No violation.")
        return nil
    }

    condition := policyConditions[0].(map[string]interface{})
    operator, _ := condition["operator"].(string)  // e.g., "IS" or "IS_NOT"
    valueStr, _ := condition["value"].(string)

    // Parse value as JSON to extract the hash
    var valObj map[string]string
    if err := json.Unmarshal([]byte(valueStr), &valObj); err != nil {
        return fmt.Errorf("failed to parse condition value: %v", err)
    }

    algorithmValue := valObj["value"]

    projects, ok := policy["projects"].([]interface{})
    if !ok || len(projects) == 0 {
        fmt.Println("No projects associated with the policy. No violation.")
        return nil
    }

    // Prepare a data structure to hold results for tabulation
    // Each record: Policy, Component, Violation State
    var results [][]string

    for _, p := range projects {
        proj := p.(map[string]interface{})
        projectUUID, _ := proj["uuid"].(string)

        // Get components for this project
        components, err := client.GetComponentsByProjectUUID(projectUUID)
        if err != nil {
            return fmt.Errorf("failed to get components for project %s: %v", projectUUID, err)
        }

        // Evaluate condition according to the corrected logic
        for _, comp := range components {
            compHash := strings.ToLower(strings.TrimSpace(comp.Sha256))
            algoValLower := strings.ToLower(strings.TrimSpace(algorithmValue))

            // Default to no violation unless conditions met
            violation := false

            switch operator {
            case "IS":
                // IS:
                // If match => violation
                // If no match => no violation
                if compHash == algoValLower {
                    violation = true
                }
            case "IS_NOT":
                // IS_NOT:
                // If match => no violation
                // If no match => violation
                if compHash != algoValLower {
                    violation = true
                }
            default:
                // If we encounter an operator not handled, treat as no violation
                // but print a warning
                fmt.Printf("Warning: Operator %s not handled, defaulting to no violation.\n", operator)
            }

            violationState := "NOT VIOLATED"
            if violation {
                violationState = "VIOLATED"
            }

            // Add the result line for each component
            results = append(results, []string{policyName, comp.Name, violationState})
        }
    }

    if len(results) == 0 {
        // If no components or nothing processed means no violation lines
        fmt.Println("No violation detected.")
        return nil
    }

    printTabulatedResults(results)

    return nil
}

func printTabulatedResults(results [][]string) {
    // Initialize a tabwriter
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

    // Print headers
    fmt.Fprintln(w, "Policy\tComponent\tViolation State")

    // Print a separator line (optional)
    fmt.Fprintln(w, "------\t---------\t--------------")

    for _, row := range results {
        fmt.Fprintf(w, "%s\t%s\t%s\n", row[0], row[1], row[2])
    }

    w.Flush()
}
