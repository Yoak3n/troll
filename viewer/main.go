package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/Yoak3n/troll-scanner/model"
	"github.com/Yoak3n/troll-scanner/model/dto"
)

func main() {
	var (
		target = flag.String("target", "", "target topic")
		user   = flag.String("user", "", "target user")
	)
	flag.Parse()
	if *target == "" {
		fmt.Println("You must specify target topic")
		return
	}
	fs, err := os.ReadDir("data")
	if err != nil {
		fmt.Println(err)
	}
	cache := make([]string, 0)
	for _, file := range fs {
		if file.IsDir() {
			cache = append(cache, file.Name())
		}
	}
	if len(cache) == 0 {
		fmt.Println("Your cache must specify at least one topic")
		return
	}

	if !slices.Contains(cache, *target) {
		fmt.Println("Your cache must specify target topic")
		return
	}
	// 有当前话题
	videos, err := os.ReadDir(filepath.Join("data", *target))
	if err != nil {
		fmt.Println(err)
		return
	}
	store := make([]dto.VideoDataOutput, 0)
	for _, video := range videos {
		if video.IsDir() {
			continue
		}

		buf, err := os.ReadFile(filepath.Join("data", *target, video.Name()))
		if err != nil {
			fmt.Println(err)
			return
		}
		data := &dto.VideoDataOutput{}
		err = json.Unmarshal(buf, data)
		if err != nil {
			fmt.Println(err)
			return
		}
		store = append(store, *data)
	}
	if *user != "" {
		getUserComment(store, *user)
	} else {
		analyze(store)
	}

}

func analyze(data []dto.VideoDataOutput) {
	userMap := make(map[string]int)
	commentMap := make(map[string]int)
	for _, d := range data {
		for _, v := range d.Comments {
			userCount, uOk := userMap[v.Author.Name]
			if uOk {
				userMap[v.Author.Name] = userCount + 1
			} else {
				userMap[v.Author.Name] = 1
			}
			commentCount, cOk := commentMap[v.Text]
			if cOk {
				commentMap[v.Text] = commentCount + 1
			} else {
				commentMap[v.Text] = 1
			}
			for _, sub := range v.Children {
				userCount, uOk = userMap[sub.Author.Name]
				if uOk {
					userMap[sub.Author.Name] = userCount + 1
				} else {
					userMap[sub.Author.Name] = 1
				}
				commentCount, cOk = commentMap[sub.Text]
				if cOk {
					commentMap[sub.Text] = commentCount + 1
				} else {
					commentMap[sub.Text] = 1
				}
			}
		}
	}
	u := getTopN(userMap, 5)
	c := getTopN(commentMap, 20)
	fmt.Println(u)
	fmt.Println(c)
}

type KeyValue struct {
	Key   string
	Value int
}

var filterComments = []string{
	"哈哈", "哈哈哈", "哈哈哈哈", "哈哈哈哈哈", "1", "2", "来了", "来了[doge]",
	"[doge]", "[吃瓜]", "[笑哭]", "转人工", "[doge_金箍]", "@AI视频小助理 总结一下", "@AI视频总结", "[打call]",
	"前排", "6", "[星星眼]", "狼来了", "第一", "捉", "+1", "便宜", "是的", "没有", "包的", "确实", "招笑",
	"可以", "第一[doge]", "牛逼", "笑死我了", "难绷", "支持", "师傅你是做什么工作的", "刚刚", "666",
	"神人", "[打call][打call][打call]", "还真是",
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
	fmt.Printf("User %s comment list:\n", name)
	for _, c := range ret {
		fmt.Println(c.Text)
	}
	return ret
}
