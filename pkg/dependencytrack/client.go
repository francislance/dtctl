package dependencytrack

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strings"
)

type Client struct {
    BaseURL    string
    APIToken   string
    HTTPClient *http.Client
}

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

// Component represents a component in Dependency-Track.
type Component struct {
    Name   string `json:"name"`
    UUID   string `json:"uuid"`
    Sha256 string `json:"sha256"`
    // Add other fields if necessary
}

// Policy represents a policy in Dependency-Track.
type Policy struct {
    Name     string    `json:"name"`
    UUID     string    `json:"uuid"`
    Projects []Project `json:"projects,omitempty"`
    // Add other fields if necessary
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
