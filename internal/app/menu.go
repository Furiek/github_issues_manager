package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Furiek/github_issues_manager/internal/config"
	"github.com/Furiek/github_issues_manager/internal/githubapi"
)

var mainMenuItems = []string{
	"Command",
	"CRUD operator helper",
	"Issues list",
	"Exit",
}

var allowedIssueStates = map[string]struct{}{
	"open":   {},
	"closed": {},
}

var allowedIssueStateReasons = map[string]struct{}{
	"completed":   {},
	"not_planned": {},
	"reopened":    {},
}

// runMenu drives the top-level interactive application menu.
func runMenu() error {
	for {
		idx, err := selectMenu("Welcome to Github Issues Manager", mainMenuItems)
		if err != nil {
			return err
		}
		if idx < 0 {
			continue
		}

		switch idx {
		case 0:
			if err := runCommandMode(); err != nil {
				printJSONError(err)
				waitForEnter("Press Enter to return to menu...")
			}
		case 1:
			if err := runCRUDHelper(); err != nil {
				printJSONError(err)
				waitForEnter("Press Enter to return to menu...")
			}
		case 2:
			exit, err := runIssuesList()
			if err != nil {
				printJSONError(err)
				waitForEnter("Press Enter to return to menu...")
			}
			if exit {
				return nil
			}
		case 3:
			return nil
		}
	}
}

// printBanner prints the application banner with ANSI color styling.
func printBanner() {
	fmt.Print("\x1b[95m")
	fmt.Print(`
 ███████╗██╗   ██╗██████╗ ██╗███████╗██╗  ██╗
 ██╔════╝██║   ██║██╔══██╗██║██╔════╝██║ ██╔╝
 █████╗  ██║   ██║██████╔╝██║█████╗  █████╔╝
 ██╔══╝  ██║   ██║██╔══██╗██║██╔══╝  ██╔═██╗
 ██║     ╚██████╔╝██║  ██║██║███████╗██║  ██╗
 ╚═╝      ╚═════╝ ╚═╝  ╚═╝╚═╝╚══════╝╚═╝  ╚═╝
   ╔══════════════════════════════════════╗
   ║     Furiek's Github Issues Manager   ║
   ╚══════════════════════════════════════╝
`)
	fmt.Print("\x1b[0m")
}

// commandRequest is the JSON shape accepted by command mode.
type commandRequest struct {
	Command string                 `json:"command"`
	Number  int                    `json:"number,omitempty"`
	Query   string                 `json:"query,omitempty"`
	Issue   *githubapi.NewIssue    `json:"issue,omitempty"`
	Update  *githubapi.IssueUpdate `json:"update,omitempty"`
}

// runCommandMode reads a single JSON command and executes it.
func runCommandMode() error {
	fmt.Println()
	fmt.Println(`Enter JSON request (single line). Example:`)
	fmt.Println(`{"command":"update","number":2,"update":{"title":"new title"}}`)

	line, err := readLine("json> ")
	if err != nil {
		return err
	}
	if strings.TrimSpace(line) == "" {
		return nil
	}

	var req commandRequest
	if err := json.Unmarshal([]byte(line), &req); err != nil {
		return fmt.Errorf("invalid JSON request: %w", err)
	}
	return executeCommand(req)
}

// runCRUDHelper shows the CRUD shortcut menu and executes the selected action.
func runCRUDHelper() error {
	idx, err := selectMenu("CRUD operator helper", []string{"Create", "Read", "Update", "Back"})
	if err != nil {
		return err
	}
	if idx < 0 {
		return nil
	}

	switch idx {
	case 0:
		return helperCreate()
	case 1:
		return helperRead()
	case 2:
		return helperUpdate()
	default:
		return nil
	}
}

// helperCreate prompts for issue details and creates a new issue.
func helperCreate() error {
	title, err := readLine("Title: ")
	if err != nil {
		return err
	}
	body, err := readLine("Body (optional): ")
	if err != nil {
		return err
	}

	req := commandRequest{
		Command: "create",
		Issue: &githubapi.NewIssue{
			Title: strings.TrimSpace(title),
			Body:  strings.TrimSpace(body),
		},
	}
	return executeCommand(req)
}

// helperRead prompts for an issue number and fetches that issue.
func helperRead() error {
	raw, err := readLine("Issue number: ")
	if err != nil {
		return err
	}
	num, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || num <= 0 {
		return fmt.Errorf("issue number must be > 0")
	}
	return executeCommand(commandRequest{
		Command: "read",
		Number:  num,
	})
}

// helperUpdate prompts for an issue number and update fields, then applies the update.
func helperUpdate() error {
	raw, err := readLine("Issue number: ")
	if err != nil {
		return err
	}
	num, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || num <= 0 {
		return fmt.Errorf("issue number must be > 0")
	}

	title, err := readLine("New title (empty to skip): ")
	if err != nil {
		return err
	}
	body, err := readLine("New body (empty to skip): ")
	if err != nil {
		return err
	}
	state, err := readLine("State (empty to skip): ")
	if err != nil {
		return err
	}
	stateReason, err := readLine("State reason (empty to skip): ")
	if err != nil {
		return err
	}

	upd := &githubapi.IssueUpdate{}
	if v := strings.TrimSpace(title); v != "" {
		upd.Title = &v
	}
	if v := strings.TrimSpace(body); v != "" {
		upd.Body = &v
	}
	if v := strings.TrimSpace(state); v != "" {
		upd.State = &v
	}
	if v := strings.TrimSpace(stateReason); v != "" {
		upd.StateReason = &v
	}

	return executeCommand(commandRequest{
		Command: "update",
		Number:  num,
		Update:  upd,
	})
}

// runIssuesList shows repository issues and lets the user open an issue detail view.
func runIssuesList() (bool, error) {
	owner, repo, err := config.RepoContextFromEnv()
	if err != nil {
		return false, err
	}

	for {
		result, err := githubapi.SearchIssues([]string{fmt.Sprintf("repo:%s/%s is:issue", owner, repo)})
		if err != nil {
			return false, err
		}
		if len(result.Items) == 0 {
			fmt.Println("No issues found.")
			waitForEnter("Press Enter to return to menu...")
			return false, nil
		}

		items := make([]string, 0, len(result.Items))
		for _, it := range result.Items {
			items = append(items, fmt.Sprintf("#%d %s", it.Number, it.Title))
		}

		idx, err := selectMenu("Issues list", items)
		if err != nil {
			return false, err
		}
		if idx < 0 {
			return false, nil
		}
		if idx >= len(result.Items) {
			continue
		}

		exitApp, err := runIssueDetail(owner, repo, result.Items[idx].Number)
		if err != nil {
			printJSONError(err)
			waitForEnter("Press Enter to continue...")
			continue
		}
		if exitApp {
			return true, nil
		}
	}
}

// runIssueDetail shows one issue and handles detail actions.
func runIssueDetail(owner, repo string, number int) (bool, error) {
	for {
		issue, err := githubapi.GetIssue(owner, repo, number)
		if err != nil {
			return false, err
		}

		title := fmt.Sprintf(
			"Issue #%d\nTitle: %s\nState: %s\nState reason: %s\nBody: %s",
			issue.Number,
			emptyFallback(issue.Title, "-"),
			emptyFallback(issue.State, "-"),
			emptyFallback(issue.StateReason, "-"),
			emptyFallback(issue.Body, "-"),
		)

		idx, err := selectMenu(title, []string{"Edit", "Back", "Exit"})
		if err != nil {
			return false, err
		}
		switch idx {
		case 0:
			if err := editIssue(owner, repo, number); err != nil {
				printJSONError(err)
				waitForEnter("Press Enter to continue...")
			}
		case 1, -1:
			return false, nil
		case 2:
			return true, nil
		}
	}
}

// editIssue prompts for a field update and sends a patch request for the issue.
func editIssue(owner, repo string, number int) error {
	fieldIdx, err := selectMenu("Select field to edit", []string{
		"Title",
		"Body",
		"State",
		"State reason",
		"Assignee",
		"Assignees (comma separated)",
		"Labels (comma separated)",
		"Type",
		"Milestone",
		"Back",
	})
	if err != nil {
		return err
	}
	if fieldIdx < 0 || fieldIdx == 9 {
		return nil
	}

	value, err := readLine("New value: ")
	if err != nil {
		return err
	}
	v := strings.TrimSpace(value)

	update := &githubapi.IssueUpdate{}
	switch fieldIdx {
	case 0:
		update.Title = &v
	case 1:
		update.Body = &v
	case 2:
		update.State = &v
	case 3:
		update.StateReason = &v
	case 4:
		update.Assignee = &v
	case 5:
		update.Assignees = splitComma(v)
	case 6:
		update.Labels = splitComma(v)
	case 7:
		update.Type = &v
	case 8:
		if n, err := strconv.Atoi(v); err == nil {
			update.Milestone = n
		} else {
			update.Milestone = v
		}
	}

	if err := validateIssueUpdate(update); err != nil {
		return err
	}

	updated, err := githubapi.UpdateIssue(owner, repo, number, update)
	if err != nil {
		return err
	}
	printJSONOK(updated)
	waitForEnter("Press Enter to continue...")
	return nil
}

// splitComma parses a comma-separated string into trimmed non-empty values.
func splitComma(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(p)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// emptyFallback returns fallback when v is blank after trimming.
func emptyFallback(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

// validateIssueUpdate validates and normalizes constrained issue update fields.
func validateIssueUpdate(update *githubapi.IssueUpdate) error {
	if update == nil {
		return errors.New("missing update payload")
	}

	if update.State != nil {
		v := strings.ToLower(strings.TrimSpace(*update.State))
		if v == "" {
			return errors.New("state cannot be empty; allowed values: open, closed")
		}
		if _, ok := allowedIssueStates[v]; !ok {
			return fmt.Errorf("invalid state %q; allowed values: open, closed", *update.State)
		}
		*update.State = v
	}

	if update.StateReason != nil {
		v := strings.ToLower(strings.TrimSpace(*update.StateReason))
		if v == "" {
			return errors.New("state_reason cannot be empty; allowed values: completed, not_planned, reopened")
		}
		if _, ok := allowedIssueStateReasons[v]; !ok {
			return fmt.Errorf("invalid state_reason %q; allowed values: completed, not_planned, reopened", *update.StateReason)
		}
		*update.StateReason = v
	}

	return nil
}

// executeCommand runs one command-mode request against the GitHub API.
func executeCommand(req commandRequest) error {
	owner, repo, err := config.RepoContextFromEnv()
	if err != nil {
		return err
	}

	cmd := strings.ToLower(strings.TrimSpace(req.Command))
	switch cmd {
	case "create":
		if req.Issue == nil {
			return errors.New("missing issue payload")
		}
		if strings.TrimSpace(req.Issue.Title) == "" {
			return errors.New("issue.title is required")
		}
		issue, err := githubapi.CreateIssue(owner, repo, req.Issue)
		if err != nil {
			return err
		}
		printJSONOK(issue)
		return nil
	case "read":
		if req.Number <= 0 {
			return errors.New("number must be > 0")
		}
		issue, err := githubapi.GetIssue(owner, repo, req.Number)
		if err != nil {
			return err
		}
		printJSONOK(issue)
		return nil
	case "update":
		if req.Number <= 0 {
			return errors.New("number must be > 0")
		}
		if err := validateIssueUpdate(req.Update); err != nil {
			return err
		}
		issue, err := githubapi.UpdateIssue(owner, repo, req.Number, req.Update)
		if err != nil {
			return err
		}
		printJSONOK(issue)
		return nil
	case "list":
		result, err := githubapi.SearchIssues([]string{fmt.Sprintf("repo:%s/%s is:issue state:open", owner, repo)})
		if err != nil {
			return err
		}
		printJSONOK(result)
		return nil
	default:
		return fmt.Errorf("unsupported command %q", req.Command)
	}
}

// readLine reads one trimmed line from stdin after applying platform input mode setup.
func readLine(prompt string) (string, error) {
	if err := ensureLineInputMode(); err != nil {
		return "", err
	}
	fmt.Print(prompt)
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// printJSONOK prints a success envelope as formatted JSON.
func printJSONOK(data any) {
	resp := map[string]any{
		"ok":   true,
		"data": data,
	}
	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}

// printJSONError prints an error envelope as formatted JSON.
func printJSONError(err error) {
	resp := map[string]any{
		"ok":    false,
		"error": err.Error(),
	}
	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}
