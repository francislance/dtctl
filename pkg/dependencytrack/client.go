package dependencytrack

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strings"
)

// Client represents a Dependency-Track API client.
type Client struct {
    BaseURL    string
    APIToken   string
    HTTPClient *http.Client
}

// NewClient initializes and returns a new Client.
func NewClient(baseURL, apiToken string) *Client {
    return &Client{
        BaseURL:    strings.TrimRight(baseURL, "/"),
        APIToken:   apiToken,
        HTTPClient: &http.Client{},
    }
}

// Project represents a project in Dependency-Track.
type Project struct {
    Name string `json:"name"`
    UUID string `json:"uuid"`
    // Remove Sha256 unless it's needed
}

// ProjectReference represents the project associated with a component.
type ProjectReference struct {
    UUID string `json:"uuid"`
}

// Component represents a component in Dependency-Track.
type Component struct {
    UUID     string           `json:"uuid"`
    Name     string           `json:"name"`
    Sha256   string           `json:"sha256"`
    Sha1     string           `json:"sha1"`
    Md5      string           `json:"md5"`
    Project  ProjectReference `json:"project"`
    // Add other fields if necessary
}

// Policy represents a policy in Dependency-Track.
type Policy struct {
    Name     string    `json:"name"`
    UUID     string    `json:"uuid"`
    Projects []Project `json:"projects,omitempty"`
    // Add other fields if necessary
}

// PolicyCondition represents a policy condition in Dependency-Track.
type PolicyCondition struct {
    Operator string `json:"operator"`
    Subject  string `json:"subject"`
    Value    string `json:"value"`
    UUID     string `json:"uuid"`
}

// GetProjects fetches all projects from the Dependency-Track server.
func (c *Client) GetProjects() ([]Project, error) {
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/project", c.BaseURL), nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("X-Api-Key", c.APIToken)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get projects: %s", resp.Status)
    }
    var projects []Project
    if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
        return nil, err
    }
    return projects, nil
}

// GetProjectsByTag fetches projects filtered by a specific tag.
func (c *Client) GetProjectsByTag(tag string) ([]Project, error) {
    // URL-encode the tag to handle special characters
    encodedTag := url.PathEscape(tag)
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/project/tag/%s", c.BaseURL, encodedTag), nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("X-Api-Key", c.APIToken)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get projects by tag: %s", resp.Status)
    }
    var projects []Project
    if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
        return nil, err
    }
    return projects, nil
}

// GetComponentsByProjectUUID fetches components for a given project UUID.
func (c *Client) GetComponentsByProjectUUID(projectUUID string) ([]Component, error) {
    endpoint := fmt.Sprintf("%s/api/v1/component/project/%s", c.BaseURL, url.PathEscape(projectUUID))
    req, err := http.NewRequest("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("X-Api-Key", c.APIToken)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get components: %s", resp.Status)
    }
    var components []Component
    if err := json.NewDecoder(resp.Body).Decode(&components); err != nil {
        return nil, err
    }
    return components, nil
}

// GetComponentByUUID fetches a single component by its UUID.
func (c *Client) GetComponentByUUID(componentUUID string) (*Component, error) {
    endpoint := fmt.Sprintf("%s/api/v1/component/%s", c.BaseURL, url.PathEscape(componentUUID))
    req, err := http.NewRequest("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("X-Api-Key", c.APIToken)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get component: %s", resp.Status)
    }
    var component Component
    if err := json.NewDecoder(resp.Body).Decode(&component); err != nil {
        return nil, err
    }
    return &component, nil
}

// UpdateComponentSHA256 updates the sha256 field of a component identified by its UUID.
// Based on the working example you found, we only need uuid, name, and sha256 in the payload.
// We do a GET first to get the current name and ensure it's exactly the same name (including spaces).
func (c *Client) UpdateComponentSHA256(componentUUID, newSHA256 string) error {
    // Fetch existing component details
    existingComponent, err := c.GetComponentByUUID(componentUUID)
    if err != nil {
        return fmt.Errorf("failed to retrieve existing component: %v", err)
    }

    // Prepare the minimal payload based on the working Postman request:
    // Just uuid, name, sha256
    payload := map[string]string{
        "uuid":   existingComponent.UUID,
        "name":   existingComponent.Name,
        "sha256": newSHA256,
    }

    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %v", err)
    }

    // Use POST to /v1/component
    endpoint := fmt.Sprintf("%s/api/v1/component", c.BaseURL)
    req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
    if err != nil {
        return fmt.Errorf("failed to create POST request: %v", err)
    }

    req.Header.Set("X-Api-Key", c.APIToken)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to perform POST request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        // Attempt to read the response body for more detailed error information
        var errorResponse map[string]interface{}
        if decodeErr := json.NewDecoder(resp.Body).Decode(&errorResponse); decodeErr != nil {
            return fmt.Errorf("failed to update component: %s", resp.Status)
        }
        // If the API returns a 'message' field, include it for more context
        if msg, exists := errorResponse["message"]; exists {
            return fmt.Errorf("failed to update component: %s, message: %v", resp.Status, msg)
        }
        return fmt.Errorf("failed to update component: %s, response: %v", resp.Status, errorResponse)
    }

    return nil
}

// GetPolicies fetches all policies from the Dependency-Track server.
func (c *Client) GetPolicies() ([]Policy, error) {
    endpoint := fmt.Sprintf("%s/api/v1/policy", c.BaseURL)
    req, err := http.NewRequest("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("X-Api-Key", c.APIToken)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get policies: %s", resp.Status)
    }
    var policies []Policy
    if err := json.NewDecoder(resp.Body).Decode(&policies); err != nil {
        return nil, err
    }
    return policies, nil
}

// UpdatePolicyCondition updates a policy's condition by sending a POST request.
func (c *Client) UpdatePolicyCondition(condition PolicyCondition) error {
    jsonPayload, err := json.Marshal(condition)
    if err != nil {
        return fmt.Errorf("failed to marshal policy condition: %v", err)
    }

    endpoint := fmt.Sprintf("%s/api/v1/policy/condition", c.BaseURL)
    req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
    if err != nil {
        return fmt.Errorf("failed to create POST request: %v", err)
    }

    req.Header.Set("X-Api-Key", c.APIToken)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to perform POST request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        var errorResponse map[string]interface{}
        if decodeErr := json.NewDecoder(resp.Body).Decode(&errorResponse); decodeErr != nil {
            return fmt.Errorf("failed to update policy condition: %s", resp.Status)
        }
        if msg, exists := errorResponse["message"]; exists {
            return fmt.Errorf("failed to update policy condition: %s, message: %v", resp.Status, msg)
        }
        return fmt.Errorf("failed to update policy condition: %s, response: %v", resp.Status, errorResponse)
    }

    return nil
}

// GetPolicyByUUID fetches a single policy by its UUID.
func (c *Client) GetPolicyByUUID(policyUUID string) (map[string]interface{}, error) {
    endpoint := fmt.Sprintf("%s/api/v1/policy/%s", c.BaseURL, url.PathEscape(policyUUID))
    req, err := http.NewRequest("GET", endpoint, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create GET request: %v", err)
    }

    req.Header.Set("X-Api-Key", c.APIToken)
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to perform GET request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        var errorResponse map[string]interface{}
        if decodeErr := json.NewDecoder(resp.Body).Decode(&errorResponse); decodeErr != nil {
            return nil, fmt.Errorf("failed to get policy: %s", resp.Status)
        }
        return nil, fmt.Errorf("failed to get policy: %s, response: %v", resp.Status, errorResponse)
    }

    var policy map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
        return nil, fmt.Errorf("failed to decode policy response: %v", err)
    }

    return policy, nil
}
