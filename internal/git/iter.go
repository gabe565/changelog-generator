package git

import (
	"errors"
	"io"
	"iter"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TagIter(repo *git.Repository) iter.Seq2[*plumbing.Reference, error] {
	return func(yield func(*plumbing.Reference, error) bool) {
		tags, err := repo.Tags()
		if err != nil {
			yield(nil, err)
			return
		}
		defer tags.Close()

		for {
			ref, err := tags.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				if !yield(ref, err) {
					return
				}
			}

			if !yield(ref, nil) {
				return
			}
		}
	}
}

func CommitIter(repo *git.Repository, opts *git.LogOptions) iter.Seq2[*object.Commit, error] {
	return func(yield func(*object.Commit, error) bool) {
		commits, err := repo.Log(opts)
		if err != nil {
			if errors.Is(err, plumbing.ErrReferenceNotFound) {
				err = ErrNoCommits
			}
			yield(nil, err)
			return
		}
		defer commits.Close()

		for {
			ref, err := commits.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				if !yield(ref, err) {
					return
				}
			}

			if !yield(ref, nil) {
				return
			}
		}
	}
}
