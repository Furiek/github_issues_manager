package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Furiek/github_issues_manager/backend/github"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("  issues create <owner> <repo> [title]")
		fmt.Println("  issues search <terms...>")
		os.Exit(1)
	}

	action := strings.ToLower(os.Args[1])
	switch action {
	case "create":
		if len(os.Args) < 4 {
			log.Fatal("usage: issues create <owner> <repo> [title]")
		}
		owner, repo := os.Args[2], os.Args[3]
		title := ""
		if len(os.Args) > 4 {
			title = strings.Join(os.Args[4:], " ")
		} else {
			fmt.Print("Title: ")
			in := bufio.NewScanner(os.Stdin)
			if !in.Scan() {
				log.Fatal("failed to read title")
			}
			title = strings.TrimSpace(in.Text())
		}
		if title == "" {
			log.Fatal("title cannot be empty")
		}

		issue, err := github.CreateIssue(owner, repo, &github.NewIssue{Title: title})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created issue #%d: %s\n%s\n", issue.Number, issue.Title, issue.HTMLURL)

	case "search":
		result, err := github.SearchIssues(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Total issues found: %d\n", result.TotalCount)
		for _, it := range result.Items {
			fmt.Printf("#%d [%s] %s (by %s)\n",
				it.Number, it.State, it.Title, it.User.Login)
		}

	default:
		log.Fatalf("unknown action: %s", action)
	}
}
