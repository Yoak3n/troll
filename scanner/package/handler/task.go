package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/troll/scanner/controller"
	util2 "github.com/Yoak3n/troll/scanner/internal/util"
	"github.com/Yoak3n/troll/scanner/model"
	"github.com/Yoak3n/troll/scanner/package/util"
)

func SearchVideoOfTopic(keyword string, page int) []model.VideoData {
	params := map[string]string{
		"search_type": "video",
		"keyword":     keyword,
		"page":        strconv.Itoa(page),
		"order":       "totalrank",
	}
	addr := util.AppendParamsToUrl(SearchUrl, params)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	accountID, cookie := accountLimiter.GetAccount(ctx)
	resBuf := util2.RequestGetWithAll(addr, cookie)
	cancel()

	response := model.SearchResponse{}
	err := json.Unmarshal(resBuf, &response)
	if err != nil {
		accountLimiter.Penalize(accountID)
		logger.Logger.Errorf("Unmarshal json fail: %v", err)
		return nil
	}
	if response.Code != 0 {
		accountLimiter.Penalize(accountID)
		logger.Logger.Errorf("Get %s failed: %s", addr, response.Message)
		return nil
	}
	accountLimiter.Reward(accountID)
	videos := make([]model.VideoData, 0)
	for _, video := range response.Data.Result {
		videos = append(videos, model.VideoData{
			Avid:        video.Aid,
			Bvid:        video.Bvid,
			Title:       util2.ExtractContentWithinTag(video.Title),
			Cover:       video.Pic,
			Description: video.Description,
			Owner: model.UserData{
				Uid:  video.Mid,
				Name: video.Author,
			},
			Review: video.Review,
		})
	}
	return videos
}

func FetchVideoInfo(bvid string, topic string) *model.VideoData {
	params := map[string]string{
		"bvid": bvid,
	}
	addr := util.AppendParamsToUrl(VideoInfoUrl, params)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	accountID, cookie := accountLimiter.GetAccount(ctx)
	resBuf := util2.RequestGetWithAll(addr, cookie)
	cancel()

	response := model.VideoInfoResponse{}
	err := json.Unmarshal(resBuf, &response)
	if err != nil {
		accountLimiter.Penalize(accountID)
		logger.Logger.Errorf("Unmarshal json fail: %v", err)
		return nil
	}
	if response.Code != 0 {
		accountLimiter.Penalize(accountID)
		logger.Logger.Errorf("Get %s failed: %s", addr, response.Message)
		return nil
	}
	accountLimiter.Reward(accountID)
	video := response.Data
	ret := &model.VideoData{
		Avid:        video.Aid,
		Bvid:        video.Bvid,
		Title:       video.Title,
		Cover:       video.Pic,
		Description: video.Description,
		Owner: model.UserData{
			Uid:  video.Owner.Mid,
			Name: video.Owner.Name,
		},
		Review: int(video.Stat.Reply),
	}
	videoRecord := model.VideoTable{
		Avid:        ret.Avid,
		Title:       ret.Title,
		Bvid:        ret.Bvid,
		Description: ret.Description,
		Owner:       ret.Owner.Uid,
		Topic:       topic,
	}
	err = controller.GlobalDatabase().AddVideoRecord(videoRecord)
	if err != nil {
		logger.Logger.Errorf("AddVideoRecord err: %v", err)
	}
	AddUserByUid(ret.Owner.Uid)
	return ret
}

func FetchVideoComments(avid uint, offset string, callback ...func(int)) ([]model.CommentData, int, string) {
	comments := make([]model.CommentData, 0)
	count := 0
	params := map[string]string{
		"oid":          strconv.FormatUint(uint64(avid), 10),
		"type":         "1",
		"mode":         "3",
		"plat":         "1",
		"web_location": "1315875",
	}
	if offset != "" {
		params["pagination_str"] = url.QueryEscape(fmt.Sprintf(`{"offset":"%s"}`, offset))
	}
	uri := util.AppendParamsToUrl(LazilyLoadUrl, params)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	accountID, cookie := accountLimiter.GetAccount(ctx)
	resBuf := util2.RequestGetWithAll(uri, cookie)
	cancel()
	response := &model.LazyCommentResponse{}
	err := json.Unmarshal(resBuf, response)
	if err != nil || response.Code != 0 {
		logger.Logger.Errorf("LazilyGetAllComments err: %v %s", err, response.Message)
		accountLimiter.Penalize(accountID)
		return nil, count, ""
	}
	accountLimiter.Reward(accountID)
	if response.Data.Cursor.IsEnd {
		logger.Logger.Printf("LazilyGetAllComments %d cursor is end", avid)
		return comments, count, response.Data.Cursor.PaginationReply.NextOffset
	}
	if len(response.Data.Replies) < 1 {
		return comments, count, response.Data.Cursor.PaginationReply.NextOffset
	}
	currentComments := extractCommentsWithCallback(response.Data.Replies, 0, callback...)
	for i, v := range currentComments {
		if v.NeedExpand && len(v.Children) > 0 {
			currentComments[i] = *getCommentSubTreeWithCallback(&v, callback...)
		}
		count += len(currentComments[i].Children)
	}
	count += len(currentComments)
	logger.Logger.Printf("LazilyGetAllComments %d count: %d", avid, count)
	comments = append(comments, currentComments...)
	offset = response.Data.Cursor.PaginationReply.NextOffset
	time.Sleep(time.Second * time.Duration(rand.Intn(3)+2))
	return comments, count, offset
}

func extractCommentsWithCallback(items []model.CommentItem, parent uint, callback ...func(int)) []model.CommentData {
	comments := make([]model.CommentData, 0)
	commentsRecords := make([]model.CommentTable, 0)
	authorsRecords := make([]model.UserTable, 0)
	if parent == 0 {
		for _, cb := range callback {
			cb(len(items))
		}
	}
	for _, v := range items {
		author := model.UserData{
			Uid:      v.Mid,
			Name:     v.Member.Uname,
			Location: v.ReplyControl.Location,
		}
		authorRecord := model.UserTable{
			Username: author.Name,
			Uid:      author.Uid,
			Avatar:   v.Member.Avatar,
			Location: author.Location,
		}
		authorsRecords = append(authorsRecords, authorRecord)
		comment := model.CommentData{
			Text:   v.Content.Message,
			Author: author,
			Rpid:   v.Rpid,
			Oid:    v.Oid,
			Like:   v.Like,
		}
		commentRecord := model.CommentTable{
			Text:      comment.Text,
			Owner:     comment.Author.Uid,
			VideoAvid: comment.Oid,
			CommentId: comment.Rpid,
		}
		// if the comment without parent, then it's a top level comment
		if parent == 0 {
			comment.Children = extractCommentsWithCallback(v.Replies, v.Rpid, callback...)
			comment.NeedExpand = v.ReplyControl.SubReplyEntryText != ""
		} else {
			commentRecord.ParentComment = parent
		}
		commentsRecords = append(commentsRecords, commentRecord)
		comments = append(comments, comment)
	}
	go func() {
		err := controller.GlobalDatabase().AddUserRecord(authorsRecords)
		if err != nil {
			logger.Logger.Errorf("AddUserRecord err: %v", err)
		}
	}()
	go func() {
		err := controller.GlobalDatabase().AddCommentRecord(commentsRecords)
		if err != nil {
			logger.Logger.Errorf("AddCommentRecord err: %v", err)
		}
	}()

	return comments
}

func getCommentSubTreeWithCallback(comment *model.CommentData, callback ...func(int)) *model.CommentData {
	page := 1
	subComments := make([]model.CommentData, 0)
	times := 0
	for {
		if times >= 5 {
			break
		}
		time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
		params := map[string]string{
			"type": "1",
			"oid":  strconv.FormatUint(uint64(comment.Oid), 10),
			"root": strconv.FormatUint(uint64(comment.Rpid), 10),
			"ps":   "20",
			"pn":   strconv.Itoa(page),
		}
		uri := util.AppendParamsToUrl(SubReplyUrl, params)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		accountID, cookie := accountLimiter.GetAccount(ctx)
		resBuf := util2.RequestGetWithAll(uri, cookie)
		cancel()
		if resBuf == nil {
			logger.Logger.Errorf("getCommentSubTree err: %v", errors.New("get sub comment response returned empty string"))
			times += 1
			accountLimiter.Penalize(accountID)
			continue
		}
		response := &model.SubCommentResponse{}
		err := json.Unmarshal(resBuf, response)
		if err != nil {
			logger.Logger.Errorf("getCommentSubTree err: %v", err)
			times += 1
			accountLimiter.Penalize(accountID)
			continue
		}
		if response.Code != 0 {
			logger.Logger.Warnf("getCommentSubTree err: %v", response.Message)
			times += 1
			accountLimiter.Penalize(accountID)
			continue
		}
		accountLimiter.Reward(accountID)
		if len(response.Data.Replies) < 1 {
			break
		}
		replies := extractComments(response.Data.Replies, comment.Rpid)
		subComments = append(subComments, replies...)
		page += 1
		time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
	}
	for _, cb := range callback {
		cb(len(subComments))
	}
	comment.Children = subComments
	return comment
}
