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
