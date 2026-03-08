package utils

import (
	"fmt"
	"strings"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionSearch Action = "search"
)

var allowedActions = map[Action]struct{}{
	ActionCreate: {},
	ActionRead:   {},
	ActionUpdate: {},
	ActionDelete: {},
	ActionSearch: {},
}

func ParseAction(s string) (Action, error) {
	a := Action(strings.ToLower(strings.TrimSpace(s)))
	if _, ok := allowedActions[a]; !ok {
		return "", fmt.Errorf("invalid action %q", s)
	}
	return a, nil
}

func SplitCommaList(raw string) []string {
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

func PrintUsage() {
	fmt.Println("GitHub Issues Manager CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  issues <action> [arguments] [flags]")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  create   Create a new issue")
	fmt.Println("  read     Read a single issue by number")
	fmt.Println("  update   Update an existing issue")
	fmt.Println("  delete   Delete/close an issue")
	fmt.Println("  search   Search issues by criteria/query terms")
	fmt.Println()
	fmt.Println("Create:")
	fmt.Println("  issues create <owner> <repo> [--title <title>] [--body <body>] [--assignee <login>] [--assignees <a,b>] [--labels <l1,l2>] [--milestone <number|name>] [--type <type>] [title]")
	fmt.Println()
	fmt.Println("Read:")
	fmt.Println("  issues read <owner> <repo> <issue_number>")
	fmt.Println()
	fmt.Println("Update:")
	fmt.Println("  issues update <owner> <repo> <issue_number> [--title <title>] [--body <body>] [--assignee <login>] [--assignees <a,b>] [--labels <l1,l2>] [--milestone <number|name>] [--type <type>]")
	fmt.Println()
	fmt.Println("Delete:")
	fmt.Println("  issues delete <owner> <repo> <issue_number>")
	fmt.Println()
	fmt.Println("Search:")
	fmt.Println("  issues search <terms...>")
	fmt.Println("  issues search repo:owner/repo is:open label:bug")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  issues create octocat hello-world --title \"Login bug\" --body \"Steps to reproduce...\" --labels bug,auth")
	fmt.Println("  issues read octocat hello-world 42")
	fmt.Println("  issues update octocat hello-world 42 --title \"Updated title\"")
	fmt.Println("  issues delete octocat hello-world 42")
	fmt.Println("  issues search repo:octocat/hello-world is:open label:bug")
}
