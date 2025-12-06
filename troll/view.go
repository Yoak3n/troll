package main

import (
	"context"
	"embed"
	"fmt"
	"sync"

	"github.com/Yoak3n/troll/viewer/service/app"
	"github.com/Yoak3n/troll/scanner/package/util"
	"github.com/urfave/cli/v3"
)

//go:embed all:dist/*
var embeddedFiles embed.FS

type viewArgs struct {
	Port int
}

func viewCommand() *cli.Command {
	V := &viewArgs{}
	return &cli.Command{
		Name:  "view",
		Usage: "open data view",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Set port",
				Aliases:     []string{"p"},
				Destination: &V.Port,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if V.Port == 0 {
				V.Port = 10420
			}
			return runRouter(V.Port)
		},
	}
}

func runRouter(port int) error {
	wg := sync.WaitGroup{}
	wg.Go(func() {
		app.InitViewCommandApp(embeddedFiles, port)
	})
	util.OpenUrlOnBrowser(fmt.Sprintf("http://localhost:%d", port))
	wg.Wait()
	return nil
}
