package app

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Furiek/github_issues_manager/internal/cli"
	"github.com/Furiek/github_issues_manager/internal/config"
	"github.com/Furiek/github_issues_manager/internal/githubapi"
)

func Run() error {
	if err := config.LoadDotEnv(".env"); err != nil {
		return err
	}

	fmt.Println("Hello. Welcome to GitHub Issues Manager.")
	fmt.Println()
	cli.PrintUsage(cli.MenuMain)

	in := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nChoose action (create/read/update/delete/search/exit): ")
		if !in.Scan() {
			if err := in.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "input error: %v\n", err)
			}
			return nil
		}

		raw := strings.TrimSpace(in.Text())
		if raw == "" {
			continue
		}

		parts := strings.Fields(raw)
		actionInput := parts[0]
		args := parts[1:]

		if strings.EqualFold(actionInput, "exit") {
			fmt.Println("Exiting.")
			return nil
		}

		action, err := cli.ParseAction(actionInput)
		if err != nil {
			fmt.Printf("Unsupported action %q.\n", actionInput)
			fmt.Println("Type one of: create, read, update, delete, search, exit")
			continue
		}

		if err := runAction(action, args, in); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("\nReturning to main menu.")
		cli.PrintUsage(cli.MenuMain)
	}
}

func runAction(action cli.Action, args []string, in *bufio.Scanner) error {
	switch action {
	case cli.ActionCreate:
		return runCreate(args, in)
	case cli.ActionRead:
		cli.PrintUsage(cli.MenuRead)
		fmt.Println("Read operation is not implemented yet.")
	case cli.ActionUpdate:
		cli.PrintUsage(cli.MenuUpdate)
		fmt.Println("Update operation is not implemented yet.")
	case cli.ActionDelete:
		cli.PrintUsage(cli.MenuDelete)
		fmt.Println("Delete operation is not implemented yet.")
	case cli.ActionSearch:
		return runSearch(args)
	}
	return nil
}

func runCreate(args []string, in *bufio.Scanner) error {
	owner, repo, err := config.RepoContextFromEnv()
	if err != nil {
		return err
	}

	createFlags := flag.NewFlagSet("create", flag.ContinueOnError)
	createFlags.SetOutput(os.Stdout)
	titleFlag := createFlags.String("title", "", "Issue title")
	bodyFlag := createFlags.String("body", "", "Issue body")
	assigneeFlag := createFlags.String("assignee", "", "Single assignee login")
	assigneesFlag := createFlags.String("assignees", "", "Comma-separated assignee logins")
	labelsFlag := createFlags.String("labels", "", "Comma-separated labels")
	milestoneFlag := createFlags.String("milestone", "", "Milestone number or name")
	typeFlag := createFlags.String("type", "", "Issue type name")
	if err := createFlags.Parse(args); err != nil {
		return err
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
		if !in.Scan() {
			return fmt.Errorf("failed to read title")
		}
		title = strings.TrimSpace(in.Text())
	}
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	newIssue := &githubapi.NewIssue{
		Title: title,
		Body:  strings.TrimSpace(*bodyFlag),
	}
	if v := strings.TrimSpace(*assigneeFlag); v != "" {
		newIssue.Assignee = &v
	}
	if v := cli.SplitCommaList(*assigneesFlag); len(v) > 0 {
		newIssue.Assignees = v
	}
	if v := cli.SplitCommaList(*labelsFlag); len(v) > 0 {
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

	issue, err := githubapi.CreateIssue(owner, repo, newIssue)
	if err != nil {
		return err
	}
	fmt.Printf("Created issue #%d: %s\n%s\n", issue.Number, issue.Title, issue.HTMLURL)
	return nil
}

func runSearch(args []string) error {
	if len(args) == 0 {
		cli.PrintUsage(cli.MenuSearch)
		return nil
	}
	result, err := githubapi.SearchIssues(args)
	if err != nil {
		return err
	}
	fmt.Printf("Total issues found: %d\n", result.TotalCount)
	for _, it := range result.Items {
		fmt.Printf("#%d [%s] %s (by %s)\n", it.Number, it.State, it.Title, it.User.Login)
	}
	return nil
}

// func runUpdate(args []string) error {
// 	// owner, repo, err := config.RepoContextFromEnv()
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// return nil
// }
