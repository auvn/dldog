package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/auvn/dldog/cmd/protomod/internal/download"
)

func main() {
	app := cli.NewApp()
	app.Name = "dldog"
	app.Commands = []*cli.Command{
		&download.Cmd,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
