package download

import (
	"github.com/urfave/cli/v2"

	"github.com/auvn/dldog/internal/resource"
	"github.com/auvn/dldog/internal/yamlcfg"
)

var _flags struct {
	Config string
}

var Cmd = cli.Command{
	Name: "download",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Required:    true,
			Destination: &_flags.Config,
		},
	},

	Action: func(ctx *cli.Context) error {
		var cfg struct {
			Skip      resource.SkipConfig
			Downloads []resource.DownloadItem
		}

		err := yamlcfg.LoadFile(_flags.Config, &cfg)
		if err != nil {
			return err
		}

		return resource.Fetch(cfg.Downloads, cfg.Skip)
	},
}
