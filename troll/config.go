package main

import (
	"context"

	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/model"
	"github.com/urfave/cli/v3"
)

type configArgs struct {
	Cookie string
	Proxy  string
}

func configCommand() *cli.Command {
	C := &configArgs{}
	controller.GlobalDatabase(TrollPath, "troll")
	return &cli.Command{
		Name:  "config",
		Usage: "Set config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cookie",
				Usage:       "Set cookie",
				Aliases:     []string{"c"},
				Destination: &C.Cookie,
				Action: func(ctx context.Context, command *cli.Command, s string) error {
					return setCookie(s)
				},
			},
			&cli.StringFlag{
				Name:        "proxy",
				Usage:       "Set proxy",
				Aliases:     []string{"p"},
				Destination: &C.Proxy,
				Action: func(ctx context.Context, command *cli.Command, s string) error {
					return setProxy(s)
				},
			},
		},
	}
}

func setCookie(cookie string) error {
	return controller.GlobalDatabase().UpdateConfiguration(&model.ConfigurationTable{
		Cookie: cookie,
	})
}

func setProxy(proxy string) error {
	return controller.GlobalDatabase().UpdateConfiguration(&model.ConfigurationTable{
		Proxy: proxy,
	})
}
