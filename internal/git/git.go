package git

import (
	"errors"
	"io"

	"github.com/gabe565/changelog-generator/internal/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func FindRepo(cmd *cobra.Command) (*git.Repository, error) {
	repoPath, err := cmd.Flags().GetString("repo")
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}

	return repo, err
}

var (
	ErrNoCommits     = errors.New("no commits found")
	ErrNoPreviousTag = errors.New("no previous tag found")
)

func FindRefs(repo *git.Repository) (*plumbing.Hash, error) {
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	defer tags.Close()

	var latest *plumbing.Reference
	var previous *plumbing.Reference
	if err := tags.ForEach(func(reference *plumbing.Reference) error {
		previous = latest
		latest = reference
		return nil
	}); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	head, err := repo.Reference(plumbing.HEAD, true)
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return nil, ErrNoCommits
		}
		return nil, err
	}
	if latest == nil || head.Hash() != latest.Hash() {
		previous = latest
		latest = head
	}

	if previous == nil {
		return nil, ErrNoPreviousTag
	}

	tag, err := repo.TagObject(previous.Hash())
	switch {
	case err == nil:
		// Tag object present
		commit, err := tag.Commit()
		if err != nil {
			return nil, err
		}
		return &commit.Hash, nil
	case errors.Is(err, plumbing.ErrObjectNotFound):
		// Not a tag object
		hash := previous.Hash()
		return &hash, nil
	default:
		return nil, err
	}
}

func WalkCommits(repo *git.Repository, conf *config.Config, previous *plumbing.Hash) error {
	commits, err := repo.Log(&git.LogOptions{})
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return ErrNoCommits
		}
		return err
	}
	defer commits.Close()

	for {
		ref, err := commits.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		if previous != nil && ref.Hash == *previous {
			break
		}

		if !conf.Filters.Match(ref) {
			continue
		}

		for _, g := range conf.Groups {
			if g.Matches(ref) {
				g.AddCommit(ref)
				break
			}
		}
	}
	return nil
}
