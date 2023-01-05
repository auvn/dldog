package blob

import (
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type SyncParams struct {
	SrcDir       string
	DstDir       string
	Patterns     []string
	SkipPatterns []string
}

type Sync struct {
	Globber Globber
	Copyer  Copyer
}

func (s *Sync) Sync(params SyncParams) error {
	var syncerr error
	for _, target := range params.Patterns {
		pattern := filepath.Join(params.SrcDir, target)
		files, err := s.Globber.Glob(pattern)
		if err != nil {
			multierr.AppendInto(&syncerr,
				errors.WithMessagef(err, "Glob(%q)", pattern))
			continue
		}

		for _, src := range files {
			relpath, _ := filepath.Rel(params.SrcDir, src)
			if s.skip(params, relpath) {
				continue
			}

			dst := filepath.Join(params.DstDir, relpath)
			if err := s.Copyer.Copy(dst, src); err != nil {
				multierr.AppendInto(&syncerr,
					errors.WithMessagef(err, "Copy(%q, %q)", dst, src))
				continue
			}
		}
	}

	if syncerr != nil {
		return syncerr
	}

	return nil
}

func (s *Sync) skip(params SyncParams, path string) bool {
	for i := range params.SkipPatterns {
		if ok, _ := filepath.Match(params.SkipPatterns[i], path); ok {
			return true
		}
	}

	return false
}

type Globber interface {
	Glob(pattern string) ([]string, error)
}

type GlobberFunc func(string) ([]string, error)

func (fn GlobberFunc) Glob(pattern string) ([]string, error) {
	return fn(pattern)
}

type Copyer interface {
	Copy(dst, src string) error
}

type CopyerFunc func(string, string) error

func (fn CopyerFunc) Copy(dst, src string) error {
	return fn(dst, src)
}
