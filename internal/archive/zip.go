package archive

import (
	"archive/zip"
	"fmt"
	"path/filepath"

	"github.com/auvn/dldog/internal/fsext"
	"github.com/pkg/errors"
)

func ZipGlob(f *zip.ReadCloser, pattern string) ([]string, error) {
	var matches []string
	for _, file := range f.File {
		if file.FileInfo().IsDir() {
			continue
		}

		ok, err := filepath.Match(pattern, file.Name)
		if err != nil {
			return nil, errors.WithMessage(err, "filepath.Match")
		}

		if !ok {
			continue
		}

		matches = append(matches, file.Name)
	}

	return matches, nil
}

func ZipExtract(arch *zip.ReadCloser, dst, src string) error {
	fmt.Println("extracting zip", src, dst)
	f, err := arch.Open(src)
	if err != nil {
		return errors.WithMessagef(err, "arch.Open(%q)", src)
	}

	return fsext.ReadToFile(dst, f)
}
