package cli

import (
	"fmt"
	"strings"
)

type Action string
type MenuChoice string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionSearch Action = "search"
)

const (
	MenuMain   MenuChoice = "main"
	MenuCreate MenuChoice = "create"
	MenuRead   MenuChoice = "read"
	MenuUpdate MenuChoice = "update"
	MenuDelete MenuChoice = "delete"
	MenuSearch MenuChoice = "search"
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
