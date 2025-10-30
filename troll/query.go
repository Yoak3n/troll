package main

import (
	"context"
	"fmt"

	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/package/handler"
	"github.com/Yoak3n/troll/scanner/package/util"
	"github.com/urfave/cli/v3"
)

type queryArgs struct {
	top   string
	count int
	user  string
}

type KeyValue struct {
	Key   string
	Value float64
}

const topUser = "user"
const topComment = "comment"

func queryCommand() *cli.Command {
	q := &queryArgs{}
	return &cli.Command{
		Name:  "query",
		Usage: "query something from troll",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			err := queryEntry(q)
			if err != nil {
				return cli.Exit(err, 6)
			}
			return nil
		},
		Flags: []cli.Flag{},
		MutuallyExclusiveFlags: []cli.MutuallyExclusiveFlags{
			{
				Required: true,
				Flags: [][]cli.Flag{
					{
						&cli.StringFlag{
							Name:        "top",
							Value:       "",
							Aliases:     []string{"t"},
							Usage:       "show top users or comments",
							Destination: &q.top,
						},
						&cli.IntFlag{
							Name:        "count",
							Value:       10,
							Usage:       "limit the count of top users or comments or words ",
							Destination: &q.count,
						},
					}, {
						&cli.StringFlag{
							Name:        "user",
							Value:       "",
							Usage:       "show the comments of a user",
							Destination: &q.user,
						},
					},
				},
			},
		},
	}
}

func queryEntry(q *queryArgs) error {
	handler.Init(TrollPath, "troll")
	if q.top != "" {
		top(q.top, q.count)
	}
	if q.user != "" {
		comments, err := controller.GlobalDatabase().QueryUserCommentsListInTopic(title, q.user)
		if err != nil {
			fmt.Printf("query user comments error: %v", err)
		}
		fmt.Printf("Comments of %s under %s:", q.user, title)
		v2c := make(map[string][]string)
		for _, v := range comments {
			bvid := util.Avid2Bvid(int64(v.VideoAvid))
			_, ok := v2c[bvid]
			if !ok {
				v2c[bvid] = make([]string, 0)
			}
			v2c[bvid] = append(v2c[bvid], v.Text)
		}
		for k, v := range v2c {
			fmt.Printf("\n\nVideo: %s\n", k)
			for _, c := range v {
				fmt.Printf("  %s\n", c)
			}
		}
	}
	return nil
}

func top(flag string, count int) {
	switch flag {
	case topUser:
		type Row struct {
			Name  string
			Uid   uint
			Count int
		}
		user, err := controller.GlobalDatabase().QueryTopNUserInTopic(title, count)
		if err != nil {
			logger.Logger.Errorf("query top user error: %v", err)
		}
		rows := make([]*Row, 0)
		for _, v := range user {
			row := &Row{
				Name:  v.Username,
				Uid:   v.Uid,
				Count: v.Count,
			}
			rows = append(rows, row)
		}
		for _, row := range rows {
			fmt.Printf("| %s \t| %d | %5d |\n", row.Name, row.Uid, row.Count)
		}
	case topComment:
		comments, err := controller.GlobalDatabase().QuerySimilarComments(title, count)
		if err != nil {
			logger.Logger.Errorf("query top comment error: %v", err)
		}
		if len(comments) == 0 {
			fmt.Println("No similar comments found")
			return
		}
		for _, v := range comments {
			fmt.Printf("| %s \t| %d | %s |\n", v.Text, v.Count, v.CommentIds)
		}

	default:
		fmt.Println("Your cache choose [user,comment,word]")
	}
}
