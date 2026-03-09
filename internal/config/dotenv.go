package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// loadDotEnvFile loads KEY=VALUE pairs from a dotenv-style file.
// Existing environment variables are left unchanged.
func loadDotEnvFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for lineNum := 1; scanner.Scan(); lineNum++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("invalid dotenv line %d in %s", lineNum, path)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if key == "" {
			return fmt.Errorf("empty dotenv key on line %d in %s", lineNum, path)
		}

		if os.Getenv(key) != "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed setting %s from %s: %w", key, path, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed reading %s: %w", path, err)
	}
	return nil
}

// LoadDotEnvAuto searches for .env from the current working directory upward,
// then from the executable directory upward.
func LoadDotEnvAuto() error {
	seen := map[string]struct{}{}
	paths := []string{}

	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, dotEnvPathsFrom(cwd)...)
	}
	if exe, err := os.Executable(); err == nil {
		paths = append(paths, dotEnvPathsFrom(filepath.Dir(exe))...)
	}

	for _, p := range paths {
		clean := filepath.Clean(p)
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		if err := loadDotEnvFile(clean); err != nil {
			return err
		}
	}

	return nil
}

// dotEnvPathsFrom returns candidate .env paths from start up to the filesystem root.
func dotEnvPathsFrom(start string) []string {
	paths := []string{}
	dir := filepath.Clean(start)
	for {
		paths = append(paths, filepath.Join(dir, ".env"))
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return paths
}
