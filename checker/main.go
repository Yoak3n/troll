package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

var cache string

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}
	cmd := &cli.Command{
		Name:    "troll",
		Version: "0.0.1",
		Usage:   "search trolls from bilibili",
		Commands: []*cli.Command{
			fetchCommand(),
			queryCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cache",
				Value:       "data/cache",
				Aliases:     []string{"C"},
				Usage:       "cache path",
				Category:    "Optional",
				Destination: &cache,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("Welcome to trolls-checker")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
	//handler.NewHandler("机圈", "华为")
	//fmt.Println("Hello Checker")
}
