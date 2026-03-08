package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	OwnerEnvVar = "GITHUB_OWNER"
	RepoEnvVar  = "GITHUB_REPO"
)

func RepoContextFromEnv() (string, string, error) {
	owner := strings.TrimSpace(os.Getenv(OwnerEnvVar))
	repo := strings.TrimSpace(os.Getenv(RepoEnvVar))
	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("%s and %s must be set in environment", OwnerEnvVar, RepoEnvVar)
	}
	return owner, repo, nil
}
