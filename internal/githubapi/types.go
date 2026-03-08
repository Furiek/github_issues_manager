package githubapi

import "time"

const IssuesURL = "https://api.github.com/search/issues"
const ReposURL = "https://api.github.com/repos"
const APITokenEnvVar = "GITHUB_API_TOKEN"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

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

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Label struct {
	Name string
}

type Milestone struct {
	Number int
	Title  string
}

type IssueType struct {
	Name string
}

type NewIssue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body,omitempty"`
	Assignee  *string  `json:"assignee,omitempty"`
	Milestone any      `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Type      *string  `json:"type,omitempty"`
}

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
