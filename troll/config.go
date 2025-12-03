package main

import (
	"context"
	"errors"

	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/model"
	"github.com/urfave/cli/v3"
	"gorm.io/gorm"
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
		Type: "cookie",
		Data: cookie,
	})
}

func setProxy(proxy string) error {
	currentProxy, err := controller.GlobalDatabase().QueryConfigurationProxy()
	if err != nil {
		return controller.GlobalDatabase().UpdateConfiguration(&model.ConfigurationTable{
			Type: "proxy",
			Data: proxy,
		})
	}
	if currentProxy != nil {
		return controller.GlobalDatabase().UpdateConfiguration(&model.ConfigurationTable{
			Model: gorm.Model{
				ID: currentProxy.ID,
			},
			Type: "proxy",
			Data: proxy,
		})
	}
	return errors.New("set proxy failed")
}
