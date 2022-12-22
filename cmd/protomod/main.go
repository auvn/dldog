package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var _flags struct {
	Import struct {
		Config string
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "protomod"
	app.Commands = []*cli.Command{
		&_import,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
