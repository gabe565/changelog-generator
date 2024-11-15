package git

import (
	"errors"
	"slices"

	"gabe565.com/changelog-generator/internal/config"
	"gabe565.com/utils/must"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

func FindRepo(cmd *cobra.Command) (*git.Repository, error) {
	repoPath := must.Must2(cmd.Flags().GetString(config.FlagRepo))
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

func FindPreviousTag(repo *git.Repository, conf *config.Config) (*plumbing.Hash, error) {
	type refCommit struct {
		ref    *plumbing.Reference
		commit *object.Commit
		hash   *plumbing.Hash
	}

	var tags []refCommit //nolint:prealloc
	for ref, err := range TagIter(repo) {
		if err != nil {
			return nil, err
		}

		if !conf.Tag.Match(ref.Name().Short()) {
			continue
		}

		hash, err := getRefHash(repo, ref)
		if err != nil {
			return nil, err
		}

		commit, err := repo.CommitObject(*hash)
		if err != nil {
			return nil, err
		}

		tags = append(tags, refCommit{ref, commit, hash})
	}

	slices.SortStableFunc(tags, func(a, b refCommit) int {
		return int(a.commit.Author.When.Sub(b.commit.Author.When))
	})

	if len(tags) == 0 {
		return nil, ErrNoPreviousTag
	}

	head, err := repo.Reference(plumbing.HEAD, true)
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return nil, ErrNoCommits
		}
		return nil, err
	}

	if head.Hash() != *tags[len(tags)-1].hash {
		return tags[len(tags)-1].hash, nil
	}
	if len(tags) == 1 {
		return nil, ErrNoPreviousTag
	}
	return tags[len(tags)-2].hash, nil
}

func WalkCommits(repo *git.Repository, conf *config.Config, previous *plumbing.Hash) error {
	for ref, err := range CommitIter(repo, &git.LogOptions{}) {
		if err != nil {
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

func getRefHash(repo *git.Repository, ref *plumbing.Reference) (*plumbing.Hash, error) {
	tag, err := repo.TagObject(ref.Hash())
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
		hash := ref.Hash()
		return &hash, nil
	default:
		return nil, err
	}
}
