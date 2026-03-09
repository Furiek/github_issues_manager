package main

import (
	"log"

	"github.com/Furiek/github_issues_manager/internal/app"
)

// main starts the application and exits on fatal startup/runtime errors.
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
