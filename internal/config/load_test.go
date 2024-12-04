package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubCmd struct {
	*cobra.Command
	tempDir string
}

func newStubCmd(t *testing.T) *stubCmd {
	cmd := &stubCmd{Command: &cobra.Command{}, tempDir: t.TempDir()}
	cmd.Flags().String(FlagConfig, "", "")
	cmd.Flags().String(FlagRepo, ".", "")
	return cmd
}

func TestLoad(t *testing.T) {
	t.Parallel()
	t.Run("no config file", func(t *testing.T) {
		cmd := newStubCmd(t)

		_, err := git.PlainInit(cmd.tempDir, false)
		require.NoError(t, err)
		require.NoError(t, cmd.Flags().Set(FlagRepo, cmd.tempDir))

		conf, err := Load(cmd.Command)
		require.NoError(t, err)
		assert.EqualValues(t, NewDefault(), conf)
		assert.Empty(t, conf.Filters.Include)
		assert.Empty(t, conf.Filters.Exclude)
		if assert.Len(t, conf.Groups, 1) {
			assert.Nil(t, conf.Groups[0].re)
		}
	})

	cfgFileTests := []struct {
		path         string
		isGoReleaser bool
	}{
		{".changelog-generator.yaml", false},
		{".changelog-generator.yml", false},
		{".goreleaser.yaml", true},
		{".goreleaser.yml", true},
	}
	for _, tt := range cfgFileTests {
		t.Run("loads config at "+tt.path, func(t *testing.T) {
			t.Parallel()
			cmd := newStubCmd(t)

			data := `filters:
  exclude:
    - "^docs"
    - "^test"
groups:
  - title: Features
    order: 0
    regexp: "^(feat)"
  - title: Fixes
    order: 1
    regexp: "^(fix|perf)"
  - title: Others
    order: 999`
			if tt.isGoReleaser {
				orig := data
				data = "changelog:\n"
				scanner := bufio.NewScanner(strings.NewReader(orig))
				for scanner.Scan() {
					data += "  " + scanner.Text() + "\n"
				}
				require.NoError(t, scanner.Err())
			}

			path := filepath.Join(cmd.tempDir, tt.path)
			require.NoError(t, cmd.Flags().Set(FlagConfig, path))
			require.NoError(t, os.WriteFile(path, []byte(data), 0o666))

			conf, err := Load(cmd.Command)
			require.NoError(t, err)
			assert.Empty(t, conf.Filters.Include)
			assert.Len(t, conf.Filters.Exclude, 2)
			assert.Len(t, conf.Groups, 3)
			for _, g := range conf.Groups {
				if g.Regexp == "" {
					assert.Nil(t, g.re)
				} else {
					assert.NotNil(t, g.re)
				}
			}
		})
	}
}
