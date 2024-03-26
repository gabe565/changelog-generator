package config

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubCmd struct {
	*cobra.Command
	prevWd   string
	tempPath string
}

func newStubCmd() *stubCmd {
	temp, err := os.MkdirTemp("", "changelog-generator-")
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(temp); err != nil {
		panic(err)
	}
	cmd := &stubCmd{Command: &cobra.Command{}, prevWd: wd, tempPath: temp}
	cmd.Flags().String("config", "", "")
	if err := cmd.ParseFlags(os.Args); err != nil {
		panic(err)
	}
	return cmd
}

func (s *stubCmd) close() {
	if err := os.Chdir(s.prevWd); err != nil {
		panic(err)
	}
	if err := os.RemoveAll(s.tempPath); err != nil {
		panic(err)
	}
}

func TestLoad(t *testing.T) {
	t.Run("no config file", func(t *testing.T) {
		cmd := newStubCmd()
		defer cmd.close()

		conf, err := Load(cmd.Command)
		require.NoError(t, err)
		assert.Equal(t, Default, conf)
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
			defer func() {
				Default = NewDefault()
			}()
			cmd := newStubCmd()
			defer cmd.close()

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
			}

			if err := os.WriteFile(tt.path, []byte(data), 0o644); !assert.NoError(t, err) {
				return
			}

			conf, err := Load(cmd.Command)
			require.NoError(t, err)
			assert.Equal(t, Default, conf)
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
