//go:build windows

package app

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	enableEchoInput            = 0x0004
	enableLineInput            = 0x0002
	enableProcessedInput       = 0x0001
	enableVirtualTerminalInput = 0x0200
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

func selectMenu(title string, items []string) (int, error) {
	h := syscall.Handle(os.Stdin.Fd())
	var oldMode uint32
	if err := getConsoleMode(h, &oldMode); err != nil {
		return 0, err
	}

	newMode := oldMode
	newMode &^= enableLineInput | enableEchoInput
	newMode |= enableProcessedInput | enableVirtualTerminalInput
	if err := setConsoleMode(h, newMode); err != nil {
		return 0, err
	}
	defer setConsoleMode(h, oldMode)

	sel := 0
	buf := make([]byte, 8)

	for {
		renderMenu(title, items, sel)

		n, err := os.Stdin.Read(buf)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			continue
		}

		b := buf[0]
		if b == '\r' || b == '\n' {
			fmt.Print("\n")
			return sel, nil
		}

		if b == 0x1b && n >= 3 && buf[1] == '[' {
			switch buf[2] {
			case 'A':
				sel--
			case 'B':
				sel++
			}
		}

		if (b == 0x00 || b == 0xe0) && n >= 2 {
			switch buf[1] {
			case 72:
				sel--
			case 80:
				sel++
			}
		}

		if sel < 0 {
			sel = len(items) - 1
		}
		if sel >= len(items) {
			sel = 0
		}
	}
}

func renderMenu(title string, items []string, selected int) {
	fmt.Print("\x1b[2J\x1b[H")
	printBanner()
	fmt.Println(title)
	fmt.Println()
	for i, it := range items {
		prefix := "  "
		if i == selected {
			prefix = "> "
		}
		fmt.Printf("%s%s\n", prefix, it)
	}
	fmt.Println()
	fmt.Println("Use Up/Down arrows and Enter.")
}

func waitForEnter(prompt string) {
	fmt.Print(prompt)
	_, _ = readLine("")
}

func getConsoleMode(h syscall.Handle, mode *uint32) error {
	r, _, e := procGetConsoleMode.Call(uintptr(h), uintptr(unsafe.Pointer(mode)))
	if r == 0 {
		if e != syscall.Errno(0) {
			return e
		}
		return syscall.EINVAL
	}
	return nil
}

func setConsoleMode(h syscall.Handle, mode uint32) error {
	r, _, e := procSetConsoleMode.Call(uintptr(h), uintptr(mode))
	if r == 0 {
		if e != syscall.Errno(0) {
			return e
		}
		return syscall.EINVAL
	}
	return nil
}
