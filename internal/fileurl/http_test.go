package fileurl

import (
	"archive/zip"
	"fmt"
	"testing"

	"github.com/auvn/dldog/internal/archive"
	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	err := Download("/tmp/myfile", "https://github.com/google/protobuf/releases/download/v3.18.0/protoc-3.18.0-osx-x86_64.zip")
	if err != nil {
		panic(err)
	}

	r, err := zip.OpenReader("/tmp/myfile")
	require.NoError(t, err)

	fmt.Println(archive.ZipGlob(r, "bin/protoc"))
}
