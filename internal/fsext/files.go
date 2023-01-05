package fsext

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var KeepTempFiles = false

func TempDir(name string) (string, func()) {
	dir, err := os.MkdirTemp("", name)
	if err != nil {
		log.Fatal(err)
	}

	rm := func() { removeTempAll(dir) }
	return dir, rm
}

func TempFile(name string) (string, func()) {
	f, err := os.CreateTemp("", name)
	if err != nil {
		log.Fatal(err)
	}

	rm := func() { removeTempFile(f.Name()) }
	return f.Name(), rm
}

func CopyFile(dst, src string) error {
	srcf, err := os.Open(src)
	if err != nil {
		return errors.WithMessage(err, "os.Open")
	}

	defer srcf.Close()

	return ReadToFile(dst, srcf)
}

func ReadToFile(dst string, f io.Reader) error {
	if err := MkdirAll(filepath.Dir(dst)); err != nil {
		return errors.WithMessage(err, "MkdirAll")
	}

	dstf, err := os.Create(dst)
	if err != nil {
		return errors.WithMessage(err, "os.Create")
	}

	defer dstf.Close()

	if _, err = io.Copy(dstf, f); err != nil {
		return errors.WithMessage(err, "io.Copy")
	}

	return nil
}

func MkdirAll(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func removeTempAll(dir string) {
	if KeepTempFiles {
		return
	}

	_ = os.RemoveAll(dir)
}

func removeTempFile(f string) {
	if KeepTempFiles {
		return
	}

	_ = os.RemoveAll(f)
}
