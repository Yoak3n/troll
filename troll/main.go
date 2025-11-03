package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sort"

	"github.com/Yoak3n/gulu/util"
	"github.com/urfave/cli/v3"
)

var cache string
var title string
var TrollPath = "troll"

func init() {
	userConfigPath, _ := os.UserConfigDir()
	TrollPath = path.Join(userConfigPath, "troll")
	_ = util.CreateDirNotExists(TrollPath)
}

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}
	cmd := &cli.Command{
		Name:    "troll",
		Version: "0.2.3",
		Usage:   "search trolls from bilibili",
		Commands: []*cli.Command{
			fetchCommand(),
			queryCommand(),
			configCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cache",
				Value:       path.Join(TrollPath, "data/cache"),
				Aliases:     []string{"C"},
				Usage:       "cache path",
				Destination: &cache,
			},
			&cli.StringFlag{
				Name:        "title",
				Value:       "",
				Usage:       "specify title as directory",
				Aliases:     []string{"T"},
				Destination: &title,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("Welcome to trolls-troll!!!")
			fmt.Println("Please use subcommands fetch and query")
			return nil
		},
	}
	sort.Sort(cli.FlagsByName(cmd.Flags))
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
