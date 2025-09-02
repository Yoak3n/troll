package main

import (
	"context"

	"github.com/Yoak3n/troll-scanner/package/handler"
	"github.com/urfave/cli/v3"
)

type fetchArgs struct {
	title string
	topic string
	AVId  int64
	BVId  string
}

func fetchCommand() *cli.Command {
	F := &fetchArgs{}
	return &cli.Command{
		Name:  "fetch",
		Usage: "fetch comments of a topic from bilibili",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.String("bvid") != "" || cmd.Int64("avid") != -1 {
				if cmd.String("title") == "" {
					return cli.Exit("You need to specify a title to save this video's data", 400)
				}
			}
			fetchEntry(cache, F)
			return nil
		},
		MutuallyExclusiveFlags: []cli.MutuallyExclusiveFlags{
			{
				Required: true,
				Flags: [][]cli.Flag{
					{
						&cli.Int64Flag{
							Name:        "avid",
							Value:       -1,
							Aliases:     []string{"a"},
							Usage:       "specify a video by avid",
							Destination: &F.AVId,
						},
					}, {
						&cli.StringFlag{
							Name:        "bvid",
							Value:       "",
							Aliases:     []string{"b"},
							Usage:       "specify a video by bvid",
							Destination: &F.BVId,
						},
					}, {
						&cli.StringFlag{
							Name:        "topic",
							Value:       "",
							Aliases:     []string{"t"},
							Usage:       "specify many video by topic name",
							Destination: &F.topic,
						},
					},
				},
			},
		},
		Flags: []cli.Flag{},
	}
}

func fetchEntry(cache string, f *fetchArgs) {
	if f.title == "" {
		f.title = f.topic
	}
	h := handler.NewHandler(cache, f.title, f.topic, f.BVId, f.AVId)
	h.Run()
}
