package fsext

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

func TempDir(name string) (string, func()) {
	dir, err := os.MkdirTemp("", name)
	if err != nil {
		log.Fatal(err)
	}

	rm := func() {
		_ = os.RemoveAll(dir)
	}

	return dir, rm
}

func CopyFile(src, dst string) error {
	srcf, err := os.Open(src)
	if err != nil {
		return errors.WithMessage(err, "os.Open")
	}

	defer srcf.Close()

	if err := MkdirAll(filepath.Dir(dst)); err != nil {
		return errors.WithMessage(err, "MkdirAll")
	}

	dstf, err := os.Create(dst)
	if err != nil {
		return errors.WithMessage(err, "os.Create")
	}

	defer dstf.Close()

	if _, err = io.Copy(dstf, srcf); err != nil {
		return errors.WithMessage(err, "io.Copy")
	}

	return nil
}

func MkdirAll(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

type Syncer struct {
	SrcDir       string
	DstDir       string
	Patterns     []string
	SkipPatterns []string
}

func (s *Syncer) Sync() error {
	var syncerr error
	for _, target := range s.Patterns {
		pattern := filepath.Join(s.SrcDir, target)
		files, err := filepath.Glob(pattern)
		if err != nil {
			multierr.Append(syncerr,
				errors.WithMessagef(err, "filepath.Glob(%q)", pattern))
			continue
		}

		for _, src := range files {
			relpath, _ := filepath.Rel(s.SrcDir, src)
			if s.skip(relpath) {
				continue
			}

			dst := filepath.Join(s.DstDir, relpath)
			if err := CopyFile(src, dst); err != nil {
				multierr.Append(syncerr,
					errors.WithMessage(err, "fsext.CopyFile"))
				continue
			}
		}
	}

	if syncerr != nil {
		return syncerr
	}

	return nil
}

func (s *Syncer) skip(path string) bool {
	for i := range s.SkipPatterns {
		if ok, _ := filepath.Match(s.SkipPatterns[i], path); ok {
			return true
		}
	}

	return false
}
