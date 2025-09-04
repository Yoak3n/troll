package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Yoak3n/troll/scanner/model"
	"github.com/Yoak3n/troll/scanner/model/dto"
	"github.com/go-ego/gse"
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
var filterSingleWord = []string{
	"的", "是", "了", "在", "我", "有", "和", "就", "不", "人", "都", "一", "一个", "上", "也", "很", "到", "说", "要", "去", "你", "会", "着", "没有", "看", "好", "自己", "这", "那", "吧", "呢", "啊", "大", "吗", "加", "中", "们", "很",
	"回复", "莫",
}

type KeyValue struct {
	Key   string
	Value float64
}

const topUser = "user"
const topComment = "comment"
const topWord = "word"

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
	fs, err := os.ReadDir(cache)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("there is no targeted cache directory")
		} else {
			return err
		}
	}
	existDirs := make([]string, 0)
	for _, file := range fs {
		if file.IsDir() {
			existDirs = append(existDirs, file.Name())
		}
	}
	if len(existDirs) == 0 {
		return errors.New("your cache must specify at least one title")
	}

	if !slices.Contains(existDirs, title) {
		fmt.Printf("Existing titles are : %s\n", strings.Join(existDirs, ","))
		return errors.New("your cache must have targeted title")
	}

	videos, err := os.ReadDir(filepath.Join(cache, title))
	if err != nil {
		return err
	}

	store := make([]dto.VideoDataOutput, 0)

	for _, video := range videos {
		if video.IsDir() {
			continue
		}

		buf, err := os.ReadFile(filepath.Join(cache, title, video.Name()))
		if err != nil {
			return err
		}
		data := &dto.VideoDataOutput{}
		err = json.Unmarshal(buf, data)
		if err != nil {
			return err
		}
		store = append(store, *data)
	}
	if q.top != "" {
		top(store, q.top, q.count)
	}
	if q.user != "" {
		getUserComment(store, q.user)
	}
	return nil
}

func top(data []dto.VideoDataOutput, flag string, count int) {
	switch flag {
	case topUser:
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
	case topComment:
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
	case topWord:
		meta, wordsCountMap := genWordCountMap(data)
		now := time.Now()
		w := getTopN(wordsCountMap, count+20)
		fmt.Println("Sorted cost:", time.Since(now))
		for _, v := range w {
			fmt.Printf("%.2f\t %s", v.Value, v.Key)
			fmt.Printf("\t\ttf:%.2f idf:%.2f count:%d\n", meta[v.Key].tf, meta[v.Key].idf, meta[v.Key].count)
		}
	default:
		fmt.Println("Your cache choose [user,comment,word]")
	}
}

type WordMetaData struct {
	tf    float64
	idf   float64
	count int
}

type TempData struct {
	text            string
	frequencyInText float64
}

func genWordCountMap(data []dto.VideoDataOutput) (map[string]WordMetaData, map[string]float64) {
	wg := sync.WaitGroup{}
	wordMap := make(map[string]int)
	wordsResult := make(map[string][]TempData)
	textChannel := make(chan string, 100)
	outputChannel := make(chan map[string]int, 5)
	outputMetaChannel := make(chan map[string][]TempData, 5)
	textCount := 0
	now := time.Now()
	for i := 0; i < 3; i++ {
		seg := &gse.Segmenter{}
		err := seg.LoadDict("zh")
		if err != nil {
			return nil, nil
		}
		err = seg.LoadDictEmbed("zh")
		if err != nil {
			return nil, nil
		}
		go func() {
			currentMap := make(map[string]int)
			currentMeta := make(map[string][]TempData)

			for text := range textChannel {
				//currentTextId := generateTextId()
				textTemp := &TempData{text: text}
				words := seg.Cut(text, true)
				words = seg.Trim(words)
				length := len(words)
				if length == 0 {
					continue
				}
				for _, word := range words {
					if slices.Contains(filterSingleWord, word) {
						continue
					}
					wordCount, ok := currentMap[word]
					if ok {
						currentMap[word] = wordCount + 1
					} else {
						currentMap[word] = 1
					}
				}
				for word, times := range currentMap {
					textTemp.frequencyInText = float64(times) / float64(len(words))
					t, ok := currentMeta[word]
					if ok {
						currentMeta[word] = append(t, *textTemp)
					} else {
						currentMeta[word] = make([]TempData, 0)
						currentMeta[word] = append(currentMeta[word], *textTemp)
					}
				}

			}
			outputMetaChannel <- currentMeta
			outputChannel <- currentMap
		}()
	}
	fmt.Println("Loaded dict cost:", time.Since(now))
	now = time.Now()
	for _, d := range data {
		for _, v := range d.Comments {
			textCount += 1
			textChannel <- v.Text
			for _, sub := range v.Children {
				textCount += 1
				textChannel <- sub.Text
			}
		}
	}
	fmt.Println("Total comments count:", textCount)
	fmt.Println("Get all comments cost:", time.Since(now))
	close(textChannel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case w := <-outputChannel:
				for k, v := range w {
					_, ok := wordMap[k]
					if ok {
						wordMap[k] += v
					} else {
						wordMap[k] = v
					}
				}
			case m := <-outputMetaChannel:
				for k, v := range m {
					_, ok := wordsResult[k]
					if ok {
						wordsResult[k] = append(wordsResult[k], v...)
					} else {
						wordsResult[k] = v
					}
				}
			case <-time.After(time.Second):
				return
			}
		}
	}()
	fmt.Println("Start to tokenize the comments...")
	now = time.Now()
	wg.Wait()
	fmt.Println("Tokenize all comments cost:", time.Since(now))
	close(outputChannel)
	close(outputMetaChannel)
	fmt.Println("Analyzed all words' data")
	now = time.Now()
	a, b := countWordsTFIDF(wordsResult, uint(textCount))
	fmt.Println("Count words' tfidf cost:", time.Since(now))
	fmt.Println("Waiting for the sorted result...")
	return a, b
}

func countWordsTFIDF(wordsResult map[string][]TempData, total uint) (map[string]WordMetaData, map[string]float64) {
	data := make(map[string]WordMetaData, len(wordsResult))
	order := make(map[string]float64, len(wordsResult))

	for word, temps := range wordsResult {
		frequency := 0.0
		textSet := make(map[string]struct{}, len(temps))

		for _, temp := range temps {
			frequency += temp.frequencyInText
			textSet[temp.text] = struct{}{}
		}
		rarity := 0.0
		globalCount := len(textSet)
		if globalCount == 0 {
			continue
		}
		rarity = float64(total) / float64(globalCount)
		tfidf := frequency * math.Log(rarity)
		order[word] = tfidf

		data[word] = WordMetaData{tf: frequency, idf: rarity, count: globalCount}

	}
	return data, order
}

func getTopN[T int | float64](m map[string]T, n int) []KeyValue {
	// 创建键值对切片
	pairs := make([]KeyValue, 0, len(m))
	for k, v := range m {
		if slices.Contains(filterComments, k) {
			continue
		}
		pairs = append(pairs, KeyValue{k, float64(v)})
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
