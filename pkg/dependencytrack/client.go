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

type Project struct {
    Name string `json:"name"`
    UUID string `json:"uuid"`
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
