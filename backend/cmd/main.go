package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Furiek/github_issues_manager/backend/github"
	"github.com/Furiek/github_issues_manager/backend/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("  issues create <owner> <repo> [--title <title>] [--body <body>] [--assignee <login>] [--assignees <a,b>] [--labels <l1,l2>] [--milestone <number|name>] [--type <type>] [title]")
		os.Exit(1)
	}
	action, err := utils.ParseAction(os.Args[1])
	if err != nil {
		log.Fatalf("Unrecognized or unsupported action was used: %s", action)
	}
	switch action {
	case "create":
		if len(os.Args) < 4 {
			log.Fatal("requires parameters: issues action <owner> <repo> [--title <title>] [title]")
		}
		owner, repo := os.Args[2], os.Args[3]

		createFlags := flag.NewFlagSet("create", flag.ContinueOnError)
		titleFlag := createFlags.String("title", "", "Issue title")
		bodyFlag := createFlags.String("body", "", "Issue body")
		assigneeFlag := createFlags.String("assignee", "", "Single assignee login")
		assigneesFlag := createFlags.String("assignees", "", "Comma-separated assignee logins")
		labelsFlag := createFlags.String("labels", "", "Comma-separated labels")
		milestoneFlag := createFlags.String("milestone", "", "Milestone number or name")
		typeFlag := createFlags.String("type", "", "Issue type name")
		if err := createFlags.Parse(os.Args[4:]); err != nil {
			log.Fatal(err)
		}

		title := strings.TrimSpace(*titleFlag)
		if title == "" {
			rest := createFlags.Args()
			if len(rest) > 0 {
				title = strings.TrimSpace(strings.Join(rest, " "))
			}
		}
		if title == "" {
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

		newIssue := &github.NewIssue{
			Title: title,
			Body:  strings.TrimSpace(*bodyFlag),
		}
		if v := strings.TrimSpace(*assigneeFlag); v != "" {
			newIssue.Assignee = &v
		}
		if v := splitCommaList(*assigneesFlag); len(v) > 0 {
			newIssue.Assignees = v
		}
		if v := splitCommaList(*labelsFlag); len(v) > 0 {
			newIssue.Labels = v
		}
		if v := strings.TrimSpace(*milestoneFlag); v != "" {
			if num, err := strconv.Atoi(v); err == nil {
				newIssue.Milestone = num
			} else {
				newIssue.Milestone = v
			}
		}
		if v := strings.TrimSpace(*typeFlag); v != "" {
			newIssue.Type = &v
		}

		issue, err := github.CreateIssue(owner, repo, newIssue)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created issue #%d: %s\n%s\n", issue.Number, issue.Title, issue.HTMLURL)

	case "update":
		if len(os.Args) < 5 {
			log.Fatal("parameters: issues create <owner> <repo> [--title <title>] [title]")
		}
		owner, repo := os.Args[2], os.Args[3]

		createFlags := flag.NewFlagSet("create", flag.ContinueOnError)
		titleFlag := createFlags.String("title", "", "Issue title")
		bodyFlag := createFlags.String("body", "", "Issue body")
		assigneeFlag := createFlags.String("assignee", "", "Single assignee login")
		assigneesFlag := createFlags.String("assignees", "", "Comma-separated assignee logins")
		labelsFlag := createFlags.String("labels", "", "Comma-separated labels")
		milestoneFlag := createFlags.String("milestone", "", "Milestone number or name")
		typeFlag := createFlags.String("type", "", "Issue type name")
		if err := createFlags.Parse(os.Args[4:]); err != nil {
			log.Fatal(err)
		}

		title := strings.TrimSpace(*titleFlag)
		if title == "" {
			rest := createFlags.Args()
			if len(rest) > 0 {
				title = strings.TrimSpace(strings.Join(rest, " "))
			}
		}
		if title == "" {
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

		newIssue := &github.NewIssue{
			Title: title,
			Body:  strings.TrimSpace(*bodyFlag),
		}
		if v := strings.TrimSpace(*assigneeFlag); v != "" {
			newIssue.Assignee = &v
		}
		if v := splitCommaList(*assigneesFlag); len(v) > 0 {
			newIssue.Assignees = v
		}
		if v := splitCommaList(*labelsFlag); len(v) > 0 {
			newIssue.Labels = v
		}
		if v := strings.TrimSpace(*milestoneFlag); v != "" {
			if num, err := strconv.Atoi(v); err == nil {
				newIssue.Milestone = num
			} else {
				newIssue.Milestone = v
			}
		}
		if v := strings.TrimSpace(*typeFlag); v != "" {
			newIssue.Type = &v
		}

		issue, err := github.CreateIssue(owner, repo, newIssue)
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

func splitCommaList(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
