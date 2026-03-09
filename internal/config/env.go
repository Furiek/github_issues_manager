package config

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	// OwnerEnvVar is the environment variable name for the repository owner.
	OwnerEnvVar = "GITHUB_OWNER"
	// RepoEnvVar is the environment variable name for the repository name.
	RepoEnvVar  = "GITHUB_REPO"
)

// RepoContextFromEnv returns owner/repo values and errors if either is missing.
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
