package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrNotAGitRepository = errors.New("not a git repository (or any of the parent directories)")
	ErrNotADirectory     = errors.New("not a directory")
)

func findGitRoot(startPath string) (string, error) {
	if stat, err := os.Stat(startPath); err != nil {
		return "", err
	} else if !stat.IsDir() {
		return "", fmt.Errorf("%w: %s", ErrNotADirectory, startPath)
	}

	path := startPath
	for {
		if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
			return path, nil
		}
		parent := filepath.Dir(path)
		if parent == path {
			return "", fmt.Errorf("%w: %s", ErrNotAGitRepository, startPath)
		}
		path = parent
	}
}
