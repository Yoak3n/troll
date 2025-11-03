package handler

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Yoak3n/gulu/logger"
	util2 "github.com/Yoak3n/gulu/util"
	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/internal/util"
	"github.com/Yoak3n/troll/scanner/model"
	"github.com/Yoak3n/troll/scanner/model/dto"
	util3 "github.com/Yoak3n/troll/scanner/package/util"
)

type Video struct {
	Bvid  string
	Avid  int64
	cache string
	name  string
}

func NewVideo(cache, name, bvid string, avid int64) *Video {
	if name == "" {
		fmt.Println("You need to specify a title to save this video's data")
		return nil
	}
	v := &Video{
		cache: cache,
		name:  name,
	}
	if avid != -1 {
		bvid = util3.Avid2Bvid(avid)
		v.Avid = avid
		v.Bvid = bvid
	} else if bvid != "" {
		if strings.HasPrefix(bvid, "http://") || strings.HasPrefix(bvid, "https://") {
			parsed, err := url.Parse(bvid)
			if err != nil {
				fmt.Println("video url parse failed", err)
				return nil
			}
			path := parsed.Path
			segments := strings.Split(strings.Trim(path, "/"), "/")
			for i, segment := range segments {
				if segment == "video" && i+1 < len(segments) {
					b := segments[i+1]
					// 验证BVID格式（BV1开头 + 10个字符）
					if strings.HasPrefix(bvid, "BV1") && len(bvid) >= 12 {
						bvid = b
						avid = util3.Bvid2Avid(bvid)
					}
				}
			}
		} else if strings.HasPrefix(bvid, "BV1") && len(bvid) >= 12 {
			avid = util3.Bvid2Avid(bvid)
		} else {
			bvid = "BV1" + bvid
		}
		v.Avid = avid
		v.Bvid = bvid
	} else {
		logger.Logger.Errorf("video %s not found", bvid)
	}
	v.Run()
	return v
}

func (v *Video) Run() {
	videoInfo := fetchVideoInfo(v)
	// TODO 或许单个视频也能使用工作池来加速
	logger.Logger.Printf("====Fetch video:《%s》 begining====", videoInfo.Title)
	videoData := NewVideoDataFromResponse(*videoInfo)
	videoRecord := model.VideoTable{
		Avid:        videoData.Avid,
		Title:       videoData.Title,
		Bvid:        videoData.Bvid,
		Description: videoData.Description,
		Owner:       videoData.Owner.Uid,
		Topic:       v.name,
	}
	err := controller.GlobalDatabase().AddVideoRecord(videoRecord)
	if err != nil {
		logger.Logger.Errorf("AddVideoRecord err: %v", err)
	}
	AddUserByUid(videoData.Owner.Uid)
	out := &dto.VideoDataOutput{
		VideoID:   v.Bvid,
		Count:     countComments(videoData),
		VideoData: videoData,
	}
	jsonData, err := json.Marshal(out)
	if err != nil {
		logger.Logger.Errorln(err)
	}
	err = util2.CreateDirNotExists(fmt.Sprintf("%s/%s", v.cache, v.name))
	if err != nil {
		return
	}
	file, err := os.Create(fmt.Sprintf("%s/%s/%s.json", v.cache, v.name, videoData.Title))
	if err != nil {
		logger.Logger.Errorln(err)
	}
	defer file.Close()
	_, err = file.Write(jsonData)
	if err != nil {
		logger.Logger.Errorln(err)
	}

	logger.Logger.Printf("====Fetch video:《%s》 completed====", videoData.Title)
}

const VideoInfoUrl = "https://api.bilibili.com/x/web-interface/wbi/view"

func fetchVideoInfo(v *Video) *model.SearchItem {
	params := map[string]string{
		"bvid": v.Bvid,
	}

	uri := util3.AppendParamsToUrl(VideoInfoUrl, params)
	resBuf := util.RequestGetWithAll(uri)
	response := &model.VideoInfoResponse{}
	err := json.Unmarshal(resBuf, response)
	if err != nil {
		logger.Logger.Errorf("fetch video info failed: %v", err)
	}
	searchItem := &model.SearchItem{
		Title:  util.ExtractContentWithinTag(response.Data.Title),
		Id:     response.Data.Aid,
		Author: response.Data.Owner.Name,
		Pic:    response.Data.Pic,
		Mid:    response.Data.Owner.Mid,
		Aid:    response.Data.Aid,
		Bvid:   response.Data.Bvid,
	}

	return searchItem
}
