package app

import (
	"github.com/Furiek/github_issues_manager/internal/config"
)

func Run() error {
	if err := config.LoadDotEnvAuto(); err != nil {
		return err
	}
	return runMenu()
}
