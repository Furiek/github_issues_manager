package config

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	OwnerEnvVar = "GITHUB_OWNER"
	RepoEnvVar  = "GITHUB_REPO"
)

func RepoContextFromEnv() (string, string, error) {
	owner := strings.TrimSpace(os.Getenv(OwnerEnvVar))
	repo := strings.TrimSpace(os.Getenv(RepoEnvVar))
	missing := []string{}
	if owner == "" {
		missing = append(missing, OwnerEnvVar)
	}
	if repo == "" {
		missing = append(missing, RepoEnvVar)
	}
	if len(missing) > 0 {
		slices.Sort(missing)
		return "", "", fmt.Errorf("missing required environment variable(s): %s", strings.Join(missing, ", "))
	}
	return owner, repo, nil
}
