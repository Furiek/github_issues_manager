//go:build !windows

package app

func ensureLineInputMode() error {
	return nil
}
