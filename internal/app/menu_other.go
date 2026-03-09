//go:build !windows

package app

import (
	"fmt"
	"strconv"
	"strings"
)

func selectMenu(title string, items []string) (int, error) {
	if len(items) == 0 {
		return -1, nil
	}

	for {
		fmt.Print("\x1b[2J\x1b[H")
		printBanner()
		fmt.Println(title)
		fmt.Println()
		for i, it := range items {
			fmt.Printf("%d. %s\n", i+1, it)
		}
		fmt.Println()
		fmt.Println("Enter a number. Press q to go back.")

		choice, err := readLine("> ")
		if err != nil {
			return 0, err
		}

		choice = strings.TrimSpace(choice)
		if strings.EqualFold(choice, "q") {
			return -1, nil
		}

		n, err := strconv.Atoi(choice)
		if err != nil || n < 1 || n > len(items) {
			continue
		}
		return n - 1, nil
	}
}

func waitForEnter(prompt string) {
	fmt.Print(prompt)
	_, _ = readLine("")
}

func ensureLineInputMode() error {
	return nil
}
