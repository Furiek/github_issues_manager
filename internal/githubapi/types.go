package githubapi

import "time"

// IssuesURL is the GitHub Search Issues endpoint.
const IssuesURL = "https://api.github.com/search/issues"
// ReposURL is the base GitHub Repositories endpoint.
const ReposURL = "https://api.github.com/repos"
// APITokenEnvVar is the environment variable name for the API token.
const APITokenEnvVar = "GITHUB_API_TOKEN"

// IssuesSearchResult represents a paged issues search response.
type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

// Issue represents a GitHub issue resource.
type Issue struct {
	Number      int
	HTMLURL     string `json:"html_url"`
	Title       string
	State       string
	StateReason string `json:"state_reason"`
	User        *User
	Assignee    *User
	Assignees   []*User
	Labels      []Label
	Milestone   *Milestone
	Type        *IssueType
	CreatedAt   time.Time `json:"created_at"`
	Body        string
}

// User represents a GitHub user in issue payloads.
type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// Label represents an issue label.
type Label struct {
	Name string
}

// Milestone represents an issue milestone.
type Milestone struct {
	Number int
	Title  string
}

// IssueType represents the issue type metadata.
type IssueType struct {
	Name string
}

// NewIssue defines fields accepted when creating an issue.
type NewIssue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body,omitempty"`
	Assignee  *string  `json:"assignee,omitempty"`
	Milestone any      `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Type      *string  `json:"type,omitempty"`
}

// IssueUpdate defines patchable fields for updating an issue.
type IssueUpdate struct {
	Title       *string  `json:"title,omitempty"`
	Body        *string  `json:"body,omitempty"`
	Assignee    *string  `json:"assignee,omitempty"`
	Milestone   any      `json:"milestone,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Assignees   []string `json:"assignees,omitempty"`
	Type        *string  `json:"type,omitempty"`
	State       *string  `json:"state,omitempty"`
	StateReason *string  `json:"state_reason,omitempty"`
}
