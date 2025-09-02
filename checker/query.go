package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/Yoak3n/troll-scanner/model"
	"github.com/Yoak3n/troll-scanner/model/dto"
	"github.com/urfave/cli/v3"
)

type queryArgs struct {
	top   string
	count int
	user  string
}

var filterComments = []string{
	"哈哈", "哈哈哈", "哈哈哈哈", "哈哈哈哈哈", "1", "2", "来了", "来了[doge]",
	"[doge]", "[吃瓜]", "[笑哭]", "[doge_金箍]", "[打call]", "[支持]",
	"转人工", "@AI视频小助理 总结一下", "@AI视频总结",
	"前排", "6", "[星星眼]", "狼来了", "第一", "捉", "+1", "便宜", "是的", "没有", "包的", "确实", "招笑",
	"可以", "第一[doge]", "牛逼", "笑死我了", "难绷", "支持", "师傅你是做什么工作的", "刚刚", "666",
	"神人", "[打call][打call][打call]", "还真是",
}

type KeyValue struct {
	Key   string
	Value int
}

const topUser = "user"
const topComment = "comment"

func queryCommand() *cli.Command {
	q := &queryArgs{}
	return &cli.Command{
		Name:  "query",
		Usage: "query something from troll",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			queryEntry(q)
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
							Usage:       "show top users or comments",
							Destination: &q.top,
						},
						&cli.IntFlag{
							Name:        "count",
							Value:       10,
							Usage:       "limit the count of top users or comments ",
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

func queryEntry(q *queryArgs) {
	fs, err := os.ReadDir(cache)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("there is no targeted cache directory")
		} else {
			fmt.Println(err)
		}
	}
	existDirs := make([]string, 0)
	for _, file := range fs {
		if file.IsDir() {
			existDirs = append(existDirs, file.Name())
		}
	}
	if len(existDirs) == 0 {
		fmt.Println("Your cache must specify at least one title")
		return
	}

	if !slices.Contains(existDirs, title) {
		fmt.Println("Your cache must have targeted title")
		fmt.Printf("Existing titles are : %s\n", strings.Join(existDirs, ","))
		return
	}

	videos, err := os.ReadDir(filepath.Join(cache, title))
	if err != nil {
		fmt.Println(err)
		return
	}

	store := make([]dto.VideoDataOutput, 0)

	for _, video := range videos {
		if video.IsDir() {
			continue
		}

		buf, err := os.ReadFile(filepath.Join(cache, title, video.Name()))
		if err != nil {
			fmt.Println(err)
			return
		}
		data := &dto.VideoDataOutput{}
		err = json.Unmarshal(buf, data)
		if err != nil {
			panic(err)
			return
		}
		store = append(store, *data)
	}
	if q.top != "" {
		top(store, q.top, q.count)
	}
	if q.user != "" {
		getUserComment(store, q.user)
	}
}

func top(data []dto.VideoDataOutput, flag string, count int) {
	if flag == topUser {
		userMap := make(map[string]int)
		for _, d := range data {
			for _, v := range d.Comments {
				userCount, ok := userMap[v.Author.Name]
				if ok {
					userMap[v.Author.Name] = userCount + 1
				} else {
					userMap[v.Author.Name] = 1
				}
				for _, sub := range v.Children {
					userCount, ok = userMap[sub.Author.Name]
					if ok {
						userMap[sub.Author.Name] = userCount + 1
					} else {
						userMap[sub.Author.Name] = 1
					}
				}
			}
		}
		u := getTopN(userMap, count)
		for _, v := range u {
			fmt.Println(v.Value, v.Key)
		}
	} else if flag == topComment {
		commentMap := make(map[string]int)
		for _, d := range data {
			for _, v := range d.Comments {
				commentCount, ok := commentMap[v.Text]
				if ok {
					commentMap[v.Text] = commentCount + 1
				} else {
					commentMap[v.Text] = 1
				}
				for _, sub := range v.Children {
					commentCount, ok = commentMap[sub.Text]
					if ok {
						commentMap[sub.Text] = commentCount + 1
					} else {
						commentMap[sub.Text] = 1
					}
				}
			}
		}
		c := getTopN(commentMap, count)
		for _, v := range c {
			fmt.Println(v.Value, v.Key)
		}
	} else {
		fmt.Println("Your cache choose [user] or [comment]")
	}

}

func getTopN(m map[string]int, n int) []KeyValue {
	// 创建键值对切片
	pairs := make([]KeyValue, 0, len(m))
	for k, v := range m {
		if slices.Contains(filterComments, k) {
			continue
		}
		pairs = append(pairs, KeyValue{k, v})
	}

	// 按值降序排序
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	// 返回前N个元素
	if n > len(pairs) {
		n = len(pairs)
	}
	return pairs[:n]
}

func getUserComment(store []dto.VideoDataOutput, name string) []model.CommentData {
	ret := make([]model.CommentData, 0)
	for _, video := range store {
		for _, comment := range video.Comments {
			if comment.Author.Name == name {
				ret = append(ret, comment)
			}
			for _, sub := range comment.Children {
				if sub.Author.Name == name {
					ret = append(ret, sub)
				}
			}
		}
	}
	fmt.Printf("User %s comments list:\n", name)
	for _, c := range ret {
		fmt.Println(c.Text)
	}
	return ret
}
