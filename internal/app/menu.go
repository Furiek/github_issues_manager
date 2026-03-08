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

func runMenu() error {
	for {
		printBanner()
		idx, err := selectMenu("Welcome to Github Issues Manager", mainMenuItems)
		if err != nil {
			return err
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
			if err := runIssuesList(); err != nil {
				printJSONError(err)
			}
			waitForEnter("Press Enter to return to menu...")
		case 3:
			return nil
		}
	}
}

func printBanner() {
	fmt.Print("\x1b[95m")
	fmt.Print(`
 ███████╗██╗   ██╗██████╗ ██╗███████╗██╗  ██╗
 ██╔════╝██║   ██║██╔══██╗██║██╔════╝██║ ██╔╝
 █████╗  ██║   ██║██████╔╝██║█████╗  █████╔╝
 ██╔══╝  ██║   ██║██╔══██╗██║██╔══╝  ██╔═██╗
 ██║     ╚██████╔╝██║  ██║██║███████╗██║  ██╗
 ╚═╝      ╚═════╝ ╚═╝  ╚═╝╚═╝╚══════╝╚═╝  ╚═╝

               furiek

   ╔══════════════════════════════════════╗
   ║     Furiek's Github Issues Manager   ║
   ╚══════════════════════════════════════╝
`)
	fmt.Print("\x1b[0m")
}

type commandRequest struct {
	Command string                 `json:"command"`
	Number  int                    `json:"number,omitempty"`
	Query   string                 `json:"query,omitempty"`
	Issue   *githubapi.NewIssue    `json:"issue,omitempty"`
	Update  *githubapi.IssueUpdate `json:"update,omitempty"`
}

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

func runCRUDHelper() error {
	idx, err := selectMenu("CRUD operator helper", []string{"Create", "Read", "Update", "Back"})
	if err != nil {
		return err
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
	state, err := readLine("State open|closed (empty to skip): ")
	if err != nil {
		return err
	}
	stateReason, err := readLine("State reason completed|not_planned|reopened (empty to skip): ")
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
		if v != "open" && v != "closed" {
			return errors.New("state must be open or closed")
		}
		upd.State = &v
	}
	if v := strings.TrimSpace(stateReason); v != "" {
		if v != "completed" && v != "not_planned" && v != "reopened" {
			return errors.New("state_reason must be completed, not_planned, or reopened")
		}
		upd.StateReason = &v
	}

	return executeCommand(commandRequest{
		Command: "update",
		Number:  num,
		Update:  upd,
	})
}

func runIssuesList() error {
	owner, repo, err := config.RepoContextFromEnv()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("repo:%s/%s is:issue state:open", owner, repo)
	result, err := githubapi.SearchIssues([]string{query})
	if err != nil {
		return err
	}
	printJSONOK(result)
	return nil
}

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
		if req.Update == nil {
			return errors.New("missing update payload")
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

func readLine(prompt string) (string, error) {
	fmt.Print(prompt)
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func printJSONOK(data any) {
	resp := map[string]any{
		"ok":   true,
		"data": data,
	}
	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}

func printJSONError(err error) {
	resp := map[string]any{
		"ok":    false,
		"error": err.Error(),
	}
	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}
