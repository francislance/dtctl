// cmd/set_hashpolicycondition.go

package cmd

import (
    "fmt"
    "dtctl/pkg/config"
    "dtctl/pkg/dependencytrack"
    "encoding/json"

    "github.com/spf13/cobra"
)

// Variables to hold flag values
var (
    hcUUID          string
    hcOperator      string
    hcSubject       string
    hcAlgorithm     string
    hcAlgorithmValue string
)

// setHashPolicyConditionCmd represents the set hashpolicycondition command
var setHashPolicyConditionCmd = &cobra.Command{
    Use:   "hashpolicycondition",
    Short: "Set or update a hash policy condition",
    RunE:  setHashPolicyCondition,
}

func init() {
    // Define flags
    setHashPolicyConditionCmd.Flags().StringVar(&hcUUID, "uuid", "", "UUID of the policy condition (required)")
    setHashPolicyConditionCmd.Flags().StringVar(&hcOperator, "operator", "", "Operator value (e.g., IS_NOT) (required)")
    setHashPolicyConditionCmd.Flags().StringVar(&hcSubject, "subject", "COMPONENT_HASH", "Subject value (default: COMPONENT_HASH)")
    setHashPolicyConditionCmd.Flags().StringVar(&hcAlgorithm, "algorithm", "", "Hash algorithm (e.g., SHA-256) (required)")
    setHashPolicyConditionCmd.Flags().StringVar(&hcAlgorithmValue, "algorithm-value", "", "Hash value (required)")

    // Mark required flags
    setHashPolicyConditionCmd.MarkFlagRequired("uuid")
    setHashPolicyConditionCmd.MarkFlagRequired("operator")
    setHashPolicyConditionCmd.MarkFlagRequired("algorithm")
    setHashPolicyConditionCmd.MarkFlagRequired("algorithm-value")

    // Add the command to the set command
    setCmd.AddCommand(setHashPolicyConditionCmd)
}

// setHashPolicyCondition handles the execution of the set hashpolicycondition command
func setHashPolicyCondition(cmd *cobra.Command, args []string) error {
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

    // Construct the value field as a JSON string
    valueObj := map[string]string{
        "algorithm": hcAlgorithm,
        "value":     hcAlgorithmValue,
    }
    valueBytes, err := json.Marshal(valueObj)
    if err != nil {
        return fmt.Errorf("failed to marshal value object: %v", err)
    }
    valueStr := string(valueBytes)

    // Create the PolicyCondition struct
    condition := dependencytrack.PolicyCondition{
        Operator: hcOperator,
        Subject:  hcSubject,
        Value:    valueStr,
        UUID:     hcUUID,
    }

    // Update the policy condition
    err = client.UpdatePolicyCondition(condition)
    if err != nil {
        return fmt.Errorf("failed to update policy condition: %v", err)
    }

    fmt.Println("Policy condition updated successfully.")
    return nil
}
