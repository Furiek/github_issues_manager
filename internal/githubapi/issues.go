package githubapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CreateIssue creates a new issue in the target repository.
func CreateIssue(owner, repo string, issue *NewIssue) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues"
	var result Issue
	if err := doJSON(http.MethodPost, url, issue, http.StatusCreated, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CloseIssue is a placeholder for closing an issue and is not implemented yet.
func CloseIssue(owner, repo string, number int) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues/" + strconv.Itoa(number)
	return nil, fmt.Errorf("not implemented: CloseIssue (%s)", url)
}

// UpdateIssue updates an existing issue by issue number.
func UpdateIssue(owner, repo string, number int, issueUpdate *IssueUpdate) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues/" + strconv.Itoa(number)
	var result Issue
	if err := doJSON(http.MethodPatch, url, issueUpdate, http.StatusOK, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssue fetches a single issue by number.
func GetIssue(owner, repo string, issueNumber int) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues/" + strconv.Itoa(issueNumber)
	var result Issue
	if err := doJSON(http.MethodGet, url, nil, http.StatusOK, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteIssue is a placeholder and currently always returns false.
func DeleteIssue(owner, repo string, issueNumber int) bool {
	_ = owner
	_ = repo
	_ = issueNumber
	return false
}

// SearchIssues queries issues using GitHub search terms.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	var result IssuesSearchResult
	if err := doJSON(http.MethodGet, IssuesURL+"?q="+q, nil, http.StatusOK, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
