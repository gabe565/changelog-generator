package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

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
	cmd := &stubCmd{Command: New("", ""), tempPath: temp}
	require.NoError(t, cmd.Flags().Set("repo", cmd.tempPath))
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

	tests := []struct {
		name             string
		createTag        bool
		createTagOptions *git.CreateTagOptions
		commits          int
		wantCommits      int
		wantErr          require.ErrorAssertionFunc
	}{
		{"no tags", false, nil, 2, 2, require.NoError},
		{"lightweight tag", true, nil, 2, 1, require.NoError},
		{"annotated tag", true, &git.CreateTagOptions{Tagger: stubAuthor(), Message: "v1.0.0"}, 2, 1, require.NoError},
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
				commits = append(commits, commitFile(t, w, cmd, "test"+strconv.Itoa(i)))
			}

			if tt.createTag {
				_, err = repo.CreateTag("v1.0.0", commits[0], tt.createTagOptions)
				require.NoError(t, err)
			}

			var buf strings.Builder
			cmd.SetOut(&buf)

			tt.wantErr(t, cmd.Execute())

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
