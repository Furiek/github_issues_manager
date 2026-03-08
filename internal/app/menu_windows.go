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
	procFlushInputBuf  = kernel32.NewProc("FlushConsoleInputBuffer")
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
	restored := false
	defer func() {
		if !restored {
			_ = setConsoleMode(h, oldMode)
		}
	}()

	sel := 0
	buf := make([]byte, 8)
	renderMenu(title, items, sel)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			continue
		}

		b := buf[0]
		if b == '\r' || b == '\n' {
			if err := setConsoleMode(h, oldMode); err != nil {
				return 0, err
			}
			restored = true
			fmt.Print("\n")
			return sel, nil
		}
		if b == 'q' || b == 'Q' {
			if err := setConsoleMode(h, oldMode); err != nil {
				return 0, err
			}
			restored = true
			return -1, nil
		}

		if b == 0x1b && n >= 3 && buf[1] == '[' {
			prev := sel
			switch buf[2] {
			case 'A':
				sel--
			case 'B':
				sel++
			}
			if sel < 0 {
				sel = len(items) - 1
			}
			if sel >= len(items) {
				sel = 0
			}
			if sel != prev {
				updateMenuSelection(items, prev, sel)
			}
			continue
		}
		if b == 0x1b && n == 1 {
			if err := setConsoleMode(h, oldMode); err != nil {
				return 0, err
			}
			restored = true
			return -1, nil
		}

		if (b == 0x00 || b == 0xe0) && n >= 2 {
			prev := sel
			switch buf[1] {
			case 72:
				sel--
			case 80:
				sel++
			}
			if sel < 0 {
				sel = len(items) - 1
			}
			if sel >= len(items) {
				sel = 0
			}
			if sel != prev {
				updateMenuSelection(items, prev, sel)
			}
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
	fmt.Println("Use Up/Down arrows and Enter. Press q to go back.")
}

func updateMenuSelection(items []string, oldSel, newSel int) {
	paintMenuItem(items, oldSel, false)
	paintMenuItem(items, newSel, true)
}

func paintMenuItem(items []string, idx int, selected bool) {
	if idx < 0 || idx >= len(items) {
		return
	}

	up := len(items) + 2 - idx
	fmt.Printf("\x1b[%dA", up)
	fmt.Print("\r\x1b[2K")
	if selected {
		fmt.Printf("> %s\n", items[idx])
	} else {
		fmt.Printf("  %s\n", items[idx])
	}
	if down := up - 1; down > 0 {
		fmt.Printf("\x1b[%dB", down)
	}
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

func ensureLineInputMode() error {
	h := syscall.Handle(os.Stdin.Fd())
	var mode uint32
	if err := getConsoleMode(h, &mode); err != nil {
		return err
	}

	mode |= enableLineInput | enableEchoInput | enableProcessedInput
	mode &^= enableVirtualTerminalInput
	if err := setConsoleMode(h, mode); err != nil {
		return err
	}
	_, _, _ = procFlushInputBuf.Call(uintptr(h))
	return nil
}
