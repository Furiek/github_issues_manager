package githubapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const githubAPIVersion = "2022-11-28"

// newGitHubRequest builds an authenticated GitHub request with required headers.
func newGitHubRequest(method, url string, body io.Reader, hasJSONBody bool) (*http.Request, error) {
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
	if hasJSONBody {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("X-GitHub-Api-Version", githubAPIVersion)
	return req, nil
}

// doJSON executes an HTTP request with an optional JSON payload and decodes JSON response.
func doJSON(method, url string, payload any, expectedStatus int, out any) error {
	var body io.Reader
	hasJSONBody := false
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(raw)
		hasJSONBody = true
	}

	req, err := newGitHubRequest(method, url, body, hasJSONBody)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	ctype := strings.ToLower(strings.TrimSpace(res.Header.Get("Content-Type")))
	if ctype == "" || !strings.Contains(ctype, "json") {
		raw, _ := io.ReadAll(res.Body)
		msg := strings.TrimSpace(string(raw))
		if msg == "" {
			msg = "empty body"
		}
		return fmt.Errorf("expected JSON response, got %q: %s", ctype, msg)
	}

	if res.StatusCode != expectedStatus {
		var apiErr struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err == nil && strings.TrimSpace(apiErr.Message) != "" {
			return fmt.Errorf("github api failed: %s (%s)", res.Status, apiErr.Message)
		}
		return fmt.Errorf("github api failed: %s", res.Status)
	}

	if out == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(out)
}
