package resource

import (
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/auvn/dldog/internal/fsext"
	"github.com/auvn/dldog/internal/gitrepo"
)

func Fetch(items []DownloadItem, skip SkipConfig) error {
	git := DefaultGit

	var errs error
	for i := range items {
		if items[i].Fetch.Git == nil {
			continue
		}

		// @TODO: repos cache
		dir, clean := fsext.TempDir("resources_" + items[i].Name)
		defer clean()

		repo := gitrepo.Desc{
			Url:    items[i].Fetch.Git.Url,
			Tag:    items[i].Fetch.Git.Tag,
			Sha:    items[i].Fetch.Git.Sha,
			Branch: items[i].Fetch.Git.Branch,
		}
		if err := git.Clone(dir, repo); err != nil {
			multierr.Append(errs,
				errors.WithMessage(err, "git.Clone()"))
			continue
		}

		syncer := fsext.Syncer{
			SrcDir:       filepath.Join(dir, items[i].Fetch.Git.Cwd),
			DstDir:       filepath.Join(items[i].Destination),
			Patterns:     items[i].Fetch.Git.Files,
			SkipPatterns: skip.Files,
		}

		if err := syncer.Sync(); err != nil {
			multierr.Append(errs,
				errors.WithMessage(err, "fsext.Sync"))
			continue
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}
