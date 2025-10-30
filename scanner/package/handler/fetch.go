package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/internal/util"
	"github.com/Yoak3n/troll/scanner/model"
	util2 "github.com/Yoak3n/troll/scanner/package/util"
)

const QueryUrl = "https://api.bilibili.com/x/v2/reply"
const SubReplyUrl = "https://api.bilibili.com/x/v2/reply/reply"
const LazilyLoadUrl = "https://api.bilibili.com/x/v2/reply/wbi/main"

func NewVideoDataFromResponse(item model.SearchItem) *model.VideoData {
	comments := LazilyGetAllComments(item.Aid)
	videoData := &model.VideoData{
		Avid:        item.Id,
		Title:       util.SanitizeFilename(item.Title),
		Bvid:        item.Bvid,
		Description: item.Description,
		Owner: model.UserData{
			Uid:  item.Mid,
			Name: item.Author,
		},
		Comments: comments,
	}
	return videoData
}

// GetAllComments risky way to crawl
func GetAllComments(avid uint) []model.CommentData {
	allComments := make([]model.CommentData, 0)
	page := 1
	for {
		params := map[string]string{
			"sort": "2",
			"oid":  strconv.FormatUint(uint64(avid), 10),
			"type": "1",
			"pn":   strconv.FormatInt(int64(page), 10),
		}
		currentPageComments, err := getComments(params)
		if err != nil {
			logger.Logger.Errorf("getComments err: %v", err)
			continue
		}
		if len(currentPageComments) < 0 {
			logger.Logger.Println("GetAllComments Fin")
			return allComments
		}
		page += 1
		// 从外部获得的评论列表，其子评论最多显示3条，需要进行展开访问
		allComments = append(allComments, currentPageComments...)
	}

}

func getComments(params map[string]string) ([]model.CommentData, error) {
	addr := util2.AppendParamsToUrl(QueryUrl, params)
	resBuf := util.RequestGetWithAll(addr)
	if resBuf == nil {
		return nil, errors.New("query response returned empty string")
	}
	response := &model.CommentResponse{}
	err := json.Unmarshal(resBuf, response)
	if err != nil {
		return nil, err
	}
	if response.Code != 0 {
		return nil, errors.New(fmt.Sprintf("Response err: %s", response.Message))
	}

	comments := extractComments(response.Data.Replies, 0)
	for i, v := range comments {
		if v.NeedExpand && len(v.Children) > 0 {
			logger.Logger.Printf("getCommentSubTree %v", v)
			comments[i] = *getCommentSubTree(&v)
		}
	}
	return comments, nil

}

func extractComments(items []model.CommentItem, parent uint) []model.CommentData {
	comments := make([]model.CommentData, 0)
	commentsRecords := make([]model.CommentTable, 0)
	authorsRecords := make([]model.UserTable, 0)
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
			comment.Children = extractComments(v.Replies, v.Rpid)
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

func getCommentSubTree(comment *model.CommentData) *model.CommentData {
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
		uri := util2.AppendParamsToUrl(SubReplyUrl, params)
		resBuf := util.RequestGetWithAll(uri)
		if resBuf == nil {
			logger.Logger.Errorf("getCommentSubTree err: %v", errors.New("get sub comment response returned empty string"))
			times += 1
			continue
		}
		response := &model.SubCommentResponse{}
		err := json.Unmarshal(resBuf, response)
		if err != nil {
			logger.Logger.Errorf("getCommentSubTree err: %v", err)
			times += 1
			continue
		}
		if response.Code != 0 {
			logger.Logger.Warnf("getCommentSubTree err: %v", response.Message)
			times += 1
			continue
		}

		if len(response.Data.Replies) < 1 {
			logger.Logger.Printf("getCommentSubTree %d completed", comment.Rpid)
			break
		}
		replies := extractComments(response.Data.Replies, comment.Rpid)
		subComments = append(subComments, replies...)
		page += 1
		time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
	}
	comment.Children = subComments
	return comment
}

func LazilyGetAllComments(avid uint) []model.CommentData {
	allComments := make([]model.CommentData, 0)
	offset := ""
	index := 1
	count := 0
	for {
		logger.Logger.Printf("LazilyGetAllComments %d index: %d", avid, index)
		params := map[string]string{
			"oid":          strconv.FormatUint(uint64(avid), 10),
			"type":         "1",
			"mode":         "3",
			"plat":         "1",
			"web_location": "1315875",
		}
		if offset != "" {
			params["pagination_str"] = url.QueryEscape(fmt.Sprintf(fmt.Sprintf(`{"offset":"%s"}`, offset)))
		}
		uri := util2.AppendParamsToUrl(LazilyLoadUrl, params)
		resBuf := util.RequestGetWithAll(uri)
		response := &model.LazyCommentResponse{}
		err := json.Unmarshal(resBuf, response)
		if err != nil || response.Code != 0 {
			logger.Logger.Errorf("LazilyGetAllComments err: %v %s", err, response.Message)
			continue
		}
		if response.Data.Cursor.IsEnd {
			logger.Logger.Printf("LazilyGetAllComments %d cursor is end", avid)
			break
		}
		if len(response.Data.Replies) < 1 {
			logger.Logger.Warnln("LazilyGetAllComments replies is empty", response.Message)
			break
		}
		currentComments := extractComments(response.Data.Replies, 0)
		for i, v := range currentComments {
			if v.NeedExpand && len(v.Children) > 0 {
				currentComments[i] = *getCommentSubTree(&v)
			}
			count += len(currentComments[i].Children)

		}
		count += len(currentComments)
		logger.Logger.Printf("LazilyGetAllComments %d count: %d", avid, count)
		allComments = append(allComments, currentComments...)
		offset = response.Data.Cursor.PaginationReply.NextOffset
		index += 1
		time.Sleep(time.Second * time.Duration(rand.Intn(3)+2))
	}
	return allComments
}
