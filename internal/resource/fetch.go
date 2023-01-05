package resource

import (
	"archive/zip"
	"log"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/auvn/dldog/internal/archive"
	"github.com/auvn/dldog/internal/blob"
	"github.com/auvn/dldog/internal/fileurl"
	"github.com/auvn/dldog/internal/fsext"
	"github.com/auvn/dldog/internal/gitrepo"
)

func Fetch(items []DownloadItem, skip SkipConfig) error {
	git := DefaultGit

	var errs error
	for i := range items {

		switch {
		case items[i].Fetch.ArchiveUrl != nil:
			f, clean := fsext.TempFile("archive_url_" + items[i].Name)
			defer clean()

			url := items[i].Fetch.ArchiveUrl.Url

			log.Printf("Downloading %q into %q...\n", url, f)

			err := fileurl.Download(f, url)
			if err != nil {
				err := errors.WithMessage(err, "fileurl.Download()")
				multierr.AppendInto(&errs, err)
				continue
			}

			format := items[i].Fetch.ArchiveUrl.Format
			if format == "" {
				format = strings.TrimPrefix(filepath.Ext(url), ".")
			}
			switch archive.Format(format) {
			case archive.FormatZip:
				arch, err := zip.OpenReader(f)
				if err != nil {
					err := errors.WithMessage(err, "zip.OpenReader")
					multierr.AppendInto(&errs, err)
					continue
				}

				globber := func(str string) ([]string, error) {
					return archive.ZipGlob(arch, str)
				}
				copyer := func(dst, src string) error {
					return archive.ZipExtract(arch, dst, src)
				}
				syncer := blob.Sync{
					Globber: blob.GlobberFunc(globber),
					Copyer:  blob.CopyerFunc(copyer),
				}

				err = syncer.Sync(blob.SyncParams{
					SrcDir:       items[i].Fetch.ArchiveUrl.Cwd,
					DstDir:       items[i].Destination,
					Patterns:     items[i].Fetch.ArchiveUrl.Files,
					SkipPatterns: skip.Files,
				})
				if err != nil {
					err := errors.WithMessage(err, "Sync")
					multierr.AppendInto(&errs, err)
					continue
				}
			default:
				err := errors.Errorf("archive url: unsupported format: %q", format)
				multierr.AppendInto(&errs, err)
				continue
			}

		case items[i].Fetch.Git != nil:
			dir, clean := fsext.TempDir("git_" + items[i].Name)
			defer clean()

			repo := gitrepo.Desc{
				Url:    items[i].Fetch.Git.Url,
				Tag:    items[i].Fetch.Git.Tag,
				Sha:    items[i].Fetch.Git.Sha,
				Branch: items[i].Fetch.Git.Branch,
			}

			log.Printf("Clonning %q into %q...\n", repo.Url, dir)
			if err := git.Clone(dir, repo); err != nil {
				multierr.AppendInto(&errs,
					errors.WithMessage(err, "git.Clone()"))
				continue
			}

			syncer := blob.Sync{
				Globber: blob.GlobberFunc(filepath.Glob),
				Copyer: blob.CopyerFunc(func(dst, src string) error {
					return fsext.CopyFile(dst, src)
				}),
			}
			sync := blob.SyncParams{
				SrcDir:       filepath.Join(dir, items[i].Fetch.Git.Cwd),
				DstDir:       filepath.Join(items[i].Destination),
				Patterns:     items[i].Fetch.Git.Files,
				SkipPatterns: skip.Files,
			}

			if err := syncer.Sync(sync); err != nil {
				multierr.AppendInto(&errs,
					errors.WithMessage(err, "Sync"))
				continue
			}
		default:
			continue
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}
