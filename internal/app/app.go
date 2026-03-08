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
	if err := config.LoadDotEnvAuto(); err != nil {
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
		return runUpdate(args)
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

func runUpdate(args []string) error {
	owner, repo, err := config.RepoContextFromEnv()
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	number := fs.Int("number", 0, "Issue number")
	title := fs.String("title", "", "Issue title")
	body := fs.String("body", "", "Issue body")
	assignee := fs.String("assignee", "", "Single assignee")
	assignees := fs.String("assignees", "", "Comma-separated assignees")
	labels := fs.String("labels", "", "Comma-separated labels")
	milestone := fs.String("milestone", "", "Milestone number or name")
	typ := fs.String("type", "", "Issue type")
	state := fs.String("state", "", "open|closed")
	stateReason := fs.String("state-reason", "", "completed|not_planned|reopened")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *number <= 0 {
		if rest := fs.Args(); len(rest) > 0 {
			n, err := strconv.Atoi(strings.TrimSpace(rest[0]))
			if err != nil || n <= 0 {
				return fmt.Errorf("number must be > 0")
			}
			*number = n
		} else {
			return fmt.Errorf("number must be > 0")
		}
	}

	u := &githubapi.IssueUpdate{}
	if v := strings.TrimSpace(*title); v != "" {
		u.Title = &v
	}
	if v := strings.TrimSpace(*body); v != "" {
		u.Body = &v
	}
	if v := strings.TrimSpace(*assignee); v != "" {
		u.Assignee = &v
	}
	if v := cli.SplitCommaList(*assignees); len(v) > 0 {
		u.Assignees = v
	}
	if v := cli.SplitCommaList(*labels); len(v) > 0 {
		u.Labels = v
	}
	if v := strings.TrimSpace(*milestone); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			u.Milestone = n
		} else {
			u.Milestone = v
		}
	}
	if v := strings.TrimSpace(*typ); v != "" {
		u.Type = &v
	}
	if v := strings.TrimSpace(*state); v != "" {
		u.State = &v
	}
	if v := strings.TrimSpace(*stateReason); v != "" {
		u.StateReason = &v
	}

	before, err := githubapi.GetIssue(owner, repo, *number)
	if err != nil {
		return err
	}

	issue, err := githubapi.UpdateIssue(owner, repo, *number, u)
	if err != nil {
		return err
	}

	fmt.Printf("Updated issue #%d: %s\n%s\n", issue.Number, issue.Title, issue.HTMLURL)
	printIssueChanges(before, issue)
	return nil
}

func printIssueChanges(before, after *githubapi.Issue) {
	fmt.Println("Changes:")
	changed := false

	changed = printFieldChange("title", before.Title, after.Title) || changed
	changed = printFieldChange("body", before.Body, after.Body) || changed
	changed = printFieldChange("state", before.State, after.State) || changed
	changed = printFieldChange("state_reason", before.StateReason, after.StateReason) || changed
	changed = printFieldChange("assignee", loginFromUser(before.Assignee), loginFromUser(after.Assignee)) || changed
	changed = printFieldChange("assignees", strings.Join(userLogins(before.Assignees), ", "), strings.Join(userLogins(after.Assignees), ", ")) || changed
	changed = printFieldChange("labels", strings.Join(labelNames(before.Labels), ", "), strings.Join(labelNames(after.Labels), ", ")) || changed
	changed = printFieldChange("milestone", milestoneText(before.Milestone), milestoneText(after.Milestone)) || changed
	changed = printFieldChange("type", issueTypeName(before.Type), issueTypeName(after.Type)) || changed

	if !changed {
		fmt.Println("  (no effective changes)")
	}
}

func printFieldChange(field, oldValue, newValue string) bool {
	if oldValue == newValue {
		return false
	}
	fmt.Printf("  %s: %q -> %q\n", field, oldValue, newValue)
	return true
}

func loginFromUser(user *githubapi.User) string {
	if user == nil {
		return ""
	}
	return user.Login
}

func userLogins(users []*githubapi.User) []string {
	out := make([]string, 0, len(users))
	for _, u := range users {
		if u == nil || strings.TrimSpace(u.Login) == "" {
			continue
		}
		out = append(out, u.Login)
	}
	return out
}

func labelNames(labels []githubapi.Label) []string {
	out := make([]string, 0, len(labels))
	for _, l := range labels {
		if strings.TrimSpace(l.Name) == "" {
			continue
		}
		out = append(out, l.Name)
	}
	return out
}

func milestoneText(m *githubapi.Milestone) string {
	if m == nil {
		return ""
	}
	return fmt.Sprintf("%d:%s", m.Number, m.Title)
}

func issueTypeName(t *githubapi.IssueType) string {
	if t == nil {
		return ""
	}
	return t.Name
}
