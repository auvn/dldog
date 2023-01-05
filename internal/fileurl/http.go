package fileurl

import (
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func Download(dest, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.WithMessagef(err, "get: %q", url)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return newResponseError(resp.Body, "non 200 status code")
	}

	f, err := os.Create(dest)
	if err != nil {
		return errors.WithMessagef(err, "os.Open(%q)", dest)
	}

	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return errors.WithMessage(err, "io.Copy")
	}

	return nil
}

func newResponseError(r io.Reader, msg string) error {
	bb, err := io.ReadAll(r)
	if err != nil {
		return errors.WithMessage(err, "read response error")
	}

	return errors.Errorf("%s: %q", msg, string(bb))
}
