package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"gabe565.com/changelog-generator/internal/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildVersion(t *testing.T) {
	t.Parallel()
	type args struct {
		version string
		commit  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"version only", args{version: "1.0.0"}, "1.0.0"},
		{"version and hash", args{version: "1.0.0", commit: "deadbeef"}, "1.0.0 (deadbeef)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, buildVersion(tt.args.version, tt.args.commit))
		})
	}
}

type stubCmd struct {
	*cobra.Command
	tempPath string
}

func newStubCmd(t *testing.T) *stubCmd {
	temp, err := os.MkdirTemp("", "changelog-generator-")
	require.NoError(t, err)
	cmd := &stubCmd{Command: New(), tempPath: temp}
	require.NoError(t, cmd.Flags().Set(config.FlagRepo, cmd.tempPath))
	cmd.SetArgs([]string{})
	return cmd
}

func (s *stubCmd) close() {
	_ = os.RemoveAll(s.tempPath)
}

func stubAuthor() *object.Signature {
	return &object.Signature{
		Name:  "Test",
		Email: "test@example.com",
		When:  time.Now(),
	}
}

func commitFile(t *testing.T, w *git.Worktree, cmd *stubCmd, name string) plumbing.Hash {
	require.NoError(t, os.WriteFile(filepath.Join(cmd.tempPath, name), []byte(name), 0o666))
	_, err := w.Add(name)
	require.NoError(t, err)

	commit, err := w.Commit(name, &git.CommitOptions{
		Author: stubAuthor(),
	})
	require.NoError(t, err)

	return commit
}

func Test_run(t *testing.T) {
	t.Parallel()

	type tag struct {
		idx  []int
		opts *git.CreateTagOptions
	}
	tests := []struct {
		name        string
		commits     int
		tag         tag
		wantCommits int
		wantErr     require.ErrorAssertionFunc
	}{
		{"no commits", 0, tag{}, 0, require.NoError},
		{"no tags", 2, tag{}, 2, require.NoError},
		{"lightweight tag as latest", 2, tag{[]int{1}, nil}, 2, require.NoError},
		{"lightweight tag as previous", 2, tag{[]int{0}, nil}, 1, require.NoError},
		{"annotated tag as latest", 2, tag{[]int{1}, &git.CreateTagOptions{Tagger: stubAuthor(), Message: "v1.0.0"}}, 2, require.NoError},
		{"annotated tag as previous", 2, tag{[]int{0}, &git.CreateTagOptions{Tagger: stubAuthor(), Message: "v1.0.0"}}, 1, require.NoError},
		{"multiple tags with latest", 3, tag{[]int{0, 2}, nil}, 2, require.NoError},
		{"multiple tags as previous", 3, tag{[]int{0, 1}, nil}, 1, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := newStubCmd(t)
			t.Cleanup(cmd.close)

			repo, err := git.PlainInit(cmd.tempPath, false)
			require.NoError(t, err)
			w, err := repo.Worktree()
			require.NoError(t, err)

			var commits []plumbing.Hash
			require.GreaterOrEqual(t, tt.commits, 0, "commits must be positive")
			for i := range tt.commits {
				commit := commitFile(t, w, cmd, "test"+strconv.Itoa(i))
				for _, tagIdx := range tt.tag.idx {
					if i == tagIdx {
						_, err = repo.CreateTag("v1.0."+strconv.Itoa(i), commit, tt.tag.opts)
						require.NoError(t, err)
					}
				}
				commits = append(commits, commit)
			}

			var buf strings.Builder
			cmd.SetOut(&buf)

			tt.wantErr(t, cmd.Execute())
			assert.Equal(t, tt.wantCommits+1, strings.Count(buf.String(), "\n"), "line count mismatch")

			want := "## Changelog\n"
			require.GreaterOrEqual(t, tt.wantCommits, 0, "wantCommits must be positive")
			for i := range tt.wantCommits {
				commitIdx := i + len(commits) - tt.wantCommits
				require.GreaterOrEqual(t, commitIdx, 0, "wantCommits must be <= commits")
				want += "- " + commits[commitIdx].String()[:8] + " " + "test" + strconv.Itoa(commitIdx) + "\n"
			}
			assert.Equal(t, want, buf.String())
		})
	}
}
