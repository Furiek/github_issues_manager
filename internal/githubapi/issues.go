package githubapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CreateIssue(owner, repo string, issue *NewIssue) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues"
	body, err := json.Marshal(issue)
	if err != nil {
		return nil, err
	}

	req, err := newGitHubRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(res.Body)

		var apiErr struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(raw, &apiErr); err == nil && strings.TrimSpace(apiErr.Message) != "" {
			return nil, fmt.Errorf("create issue failed: %s (%s)", res.Status, apiErr.Message)
		}

		msg := strings.TrimSpace(string(raw))
		if msg != "" {
			return nil, fmt.Errorf("create issue failed: %s (%s)", res.Status, msg)
		}
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

func UpdateIssue(owner, repo string, number int, issueUpdate *IssueUpdate) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues/" + strconv.Itoa(number)
	body, err := json.Marshal(issueUpdate)
	if err != nil {
		return nil, err
	}

	req, err := newGitHubRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(res.Body)

		var apiErr struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(raw, &apiErr); err == nil && strings.TrimSpace(apiErr.Message) != "" {
			return nil, fmt.Errorf("update issue failed: %s (%s)", res.Status, apiErr.Message)
		}

		msg := strings.TrimSpace(string(raw))
		if msg != "" {
			return nil, fmt.Errorf("update issue failed: %s (%s)", res.Status, msg)
		}
		return nil, fmt.Errorf("update issue failed: %s", res.Status)
	}

	var result Issue
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetIssue(owner, repo string, issueNumber int) (*Issue, error) {
	url := ReposURL + "/" + owner + "/" + repo + "/issues/" + strconv.Itoa(issueNumber)
	req, err := newGitHubRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(res.Body)

		var apiErr struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(raw, &apiErr); err == nil && strings.TrimSpace(apiErr.Message) != "" {
			return nil, fmt.Errorf("get issue failed: %s (%s)", res.Status, apiErr.Message)
		}

		msg := strings.TrimSpace(string(raw))
		if msg != "" {
			return nil, fmt.Errorf("get issue failed: %s (%s)", res.Status, msg)
		}
		return nil, fmt.Errorf("get issue failed: %s", res.Status)
	}

	var result Issue
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func DeleteIssue(owner, repo string, issueNumber int) bool {
	_ = owner
	_ = repo
	_ = issueNumber
	return false
}

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
