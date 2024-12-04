package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_findGitRoot(t *testing.T) {
	t.Parallel()
	t.Run("no git repo", func(t *testing.T) {
		t.Parallel()
		temp := t.TempDir()
		got, err := findGitRoot(temp)
		require.Error(t, err)
		assert.Equal(t, "", got)
	})

	t.Run("wd is git repo", func(t *testing.T) {
		t.Parallel()
		temp := t.TempDir()
		_, err := git.PlainInit(temp, false)
		require.NoError(t, err)
		got, err := findGitRoot(temp)
		require.NoError(t, err)
		assert.Equal(t, temp, got)
	})

	t.Run("traverse", func(t *testing.T) {
		t.Parallel()
		temp := t.TempDir()
		_, err := git.PlainInit(temp, false)
		require.NoError(t, err)
		subPath := filepath.Join(temp, "a", "b", "c")
		require.NoError(t, os.MkdirAll(subPath, 0o700))
		got, err := findGitRoot(subPath)
		require.NoError(t, err)
		assert.Equal(t, temp, got)
	})

	t.Run("wd not exist", func(t *testing.T) {
		t.Parallel()
		temp := t.TempDir()
		require.NoError(t, os.Remove(temp))
		got, err := findGitRoot(temp)
		require.Error(t, err)
		assert.Equal(t, "", got)
	})

	t.Run("wd not a directory", func(t *testing.T) {
		t.Parallel()
		temp := t.TempDir()
		_, err := git.PlainInit(temp, false)
		require.NoError(t, err)
		tempFile := filepath.Join(temp, "test.txt")
		require.NoError(t, os.WriteFile(tempFile, []byte("abc"), 0o600))
		got, err := findGitRoot(tempFile)
		require.Error(t, err)
		assert.Equal(t, "", got)
	})
}
