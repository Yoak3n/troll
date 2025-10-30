package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Yoak3n/gulu/logger"
	util2 "github.com/Yoak3n/gulu/util"
	"github.com/Yoak3n/troll/scanner/internal/util"
	"github.com/Yoak3n/troll/scanner/model"
	"github.com/Yoak3n/troll/scanner/model/dto"
	util3 "github.com/Yoak3n/troll/scanner/package/util"
)

const SearchUrl = "https://api.bilibili.com/x/web-interface/wbi/search/type"

type Topic struct {
	Name    string
	KeyWord []string
	Videos  []model.VideoData
	cache   string
	wg      sync.WaitGroup
	jobs    chan model.SearchItem
}

func NewTopic(cache string, name string, topic []string) *Topic {
	now := time.Now()
	t := &Topic{
		Name:    name,
		KeyWord: topic,
		cache:   cache,
		wg:      sync.WaitGroup{},
		jobs:    make(chan model.SearchItem),
	}
	t.fetchVideos()
	cost := time.Since(now)
	logger.Logger.Printf("%s cost %vmin", t.Name, cost.Minutes())
	return t
}

func (t *Topic) fetchVideos() {
	kw := strings.Join(t.KeyWord, ",")
	params := map[string]string{
		"search_type": "video",
		"keyword":     kw,
		"page":        "1",
		"order":       "totalrank",
	}
	addr := util3.AppendParamsToUrl(SearchUrl, params)
	resBuf := util.RequestGetWithAll(addr)
	videos := make([]model.VideoData, 0)
	response := &model.SearchResponse{}
	err := json.Unmarshal(resBuf, response)
	if err != nil {
		logger.Logger.Errorf("Unmarshal json fail: %v", err)
		return
	}
	if response.Code != 0 {
		logger.Logger.Errorf("Get %s failed: %s", addr, response.Message)
		return
	}
	for i := 1; i <= 2; i++ {
		t.wg.Add(1)
		go t.worker(i, t.jobs, &t.wg)
	}
	// Begin

	for _, value := range response.Data.Result {
		logger.Logger.Printf("====Fetch video:《%s》 begining====", value.Title)
		t.jobs <- value
		//v := NewVideoDataFromResponse(value)
		//out := &dto.VideoDataOutput{
		//	VideoID:   v.Bvid,
		//	Count:     countComments(v),
		//	VideoData: v,
		//}
		//videos = append(videos, *v)
		//jsonData, err := json.Marshal(out)
		//if err != nil {
		//	logger.Logger.Errorln(err)
		//}
		//err = util2.CreateDirNotExists(fmt.Sprintf("data/%s", t.Name))
		//if err != nil {
		//	return
		//}
		//file, err := os.Create(fmt.Sprintf("data/%s/%s.json", t.Name, v.Title))
		//if err != nil {
		//	logger.Logger.Errorln(err)
		//	continue
		//}
		//_, err = file.Write(jsonData)
		//if err != nil {
		//	logger.Logger.Errorln(err)
		//}
		//file.Close()
	}
	close(t.jobs)
	t.wg.Wait()
	t.Videos = videos
}

func (t *Topic) worker(id int, jobs <-chan model.SearchItem, wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range jobs {
		value.Title = util.ExtractContentWithinTag(value.Title)
		logger.Logger.Printf("====Woker %d Fetch video:《%s》 begining====", id, value.Title)
		v := NewVideoDataFromResponse(value)
		out := &dto.VideoDataOutput{
			VideoID:   v.Bvid,
			Count:     countComments(v),
			VideoData: v,
		}
		jsonData, err := json.Marshal(out)
		if err != nil {
			logger.Logger.Errorln(err)
		}
		err = util2.CreateDirNotExists(fmt.Sprintf("%s/%s", t.cache, t.Name))
		if err != nil {
			return
		}
		file, err := os.Create(fmt.Sprintf("%s/%s/%s.json", t.cache, t.Name, v.Title))
		if err != nil {
			logger.Logger.Errorln(err)
			continue
		}
		_, err = file.Write(jsonData)
		if err != nil {
			logger.Logger.Errorln(err)
		}
		file.Close()
		logger.Logger.Printf("====Woker %d Fetch video:《%s》 completed====", id, value.Title)
	}
}

func countComments(video *model.VideoData) uint {
	var count uint
	if len(video.Comments) > 0 {
		for _, comment := range video.Comments {
			count += 1
			l := len(comment.Children)
			count += uint(l)
		}

	}
	return count
}
