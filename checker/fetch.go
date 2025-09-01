package main

import (
	"context"
	"fmt"

	"github.com/Yoak3n/troll-scanner/package/handler"
	"github.com/urfave/cli/v3"
)

type fetchArgs struct {
	Title string
	Topic string
}

func fetchCommand() *cli.Command {
	F := &fetchArgs{}
	return &cli.Command{
		Name:  "fetch",
		Usage: "fetch comments of a topic from bilibili",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("fetch comments: ")
			fetchEntry(cache, F)
			return nil
		},
		Flags: []cli.Flag{
			//&cli.StringFlag{
			//	Name:        "cache",
			//	Value:       "data/cache",
			//	Aliases:     []string{"C"},
			//	Usage:       "cache path",
			//	Category:    "Optional",
			//	Destination: &cache,
			//},
			&cli.StringFlag{
				Name:        "topic",
				Value:       "",
				Aliases:     []string{"T"},
				Usage:       "topic name",
				Category:    "Necessary",
				Destination: &F.Topic,
			},
			&cli.StringFlag{
				Name:        "title",
				Value:       "",
				Usage:       "specify title as directory",
				Category:    "Optional",
				Destination: &F.Title,
			},
		},
	}
}

func fetchEntry(cache string, f *fetchArgs) {
	if f.Title == "" {
		f.Title = f.Topic
	}
	h := handler.NewHandler(cache, f.Title, f.Topic)
	h.Run()
}
