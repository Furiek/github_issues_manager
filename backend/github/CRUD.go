package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// TODO: these have fairly similar sending procedures, can be grouped

func CreateIssue(owner, repo string, issue *NewIssue) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues"
	body, err := json.Marshal(issue)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	token := os.Getenv(APITokenEnvVar)
	if token == "" {
		return nil, fmt.Errorf("%s is not set", APITokenEnvVar)
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create issue failed: %s", res.Status)
	}

	var result Issue
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func CloseIssue(owner, repo string, number int) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues/" + strconv.Itoa(number)
	return nil, fmt.Errorf("not implemented: CloseIssue (%s)", url)
}

func UpdateIssue(owner, repo string, number int) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues/" + strconv.Itoa(number)
	return nil, fmt.Errorf("not implemented: UpdateIssue (%s)", url)
}

func GetIssue(owner, repo string, issueNumber int) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues/" + strconv.Itoa(issueNumber)
	return nil, fmt.Errorf("not implemented: GetIssue (%s)", url)
}

func DeleteIssue(owner, repo string, issueNumber int) bool {
	_ = owner
	_ = repo
	_ = issueNumber
	return false
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
