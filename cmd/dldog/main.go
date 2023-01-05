package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/auvn/dldog/cmd/dldog/internal/download"
	"github.com/auvn/dldog/internal/fsext"
)

func main() {
	app := cli.NewApp()
	app.Name = "dldog"
	app.Commands = []*cli.Command{
		&download.Cmd,
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "keep-temp-files",
			Required:    false,
			Destination: &fsext.KeepTempFiles,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
