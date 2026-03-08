package main

import (
	"log"

	"github.com/Furiek/github_issues_manager/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
