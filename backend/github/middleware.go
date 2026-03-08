package github

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const githubAPIVersion = "2022-11-28"

func newGitHubRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	token := os.Getenv(APITokenEnvVar)
	if token == "" {
		return nil, fmt.Errorf("%s is not set", APITokenEnvVar)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", githubAPIVersion)
	return req, nil
}
