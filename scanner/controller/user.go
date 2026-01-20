package controller

import (
	"fmt"

	"github.com/Yoak3n/troll/scanner/model"
)

func (d *Database) QueryTopNUserInTopic(topic string, n int) ([]model.UserQuery, error) {
	var ret []model.UserQuery
	query := `
	SELECT
		u.*,
		COUNT(c.comment_id) AS count
	FROM
		user_tables u
		INNER JOIN
		comment_tables c ON u.uid = c.owner
		INNER JOIN
		video_tables v ON c.video_avid = v.avid
	WHERE
		v.topic = ?
	GROUP BY
		u.uid
	ORDER BY
	    count DESC
	LIMIT ?;`
	err := d.db.Raw(query, topic, n).Scan(&ret).Error
	return ret, err
}

type CommentsWithVideoQuery struct {
	model.CommentTable
	model.VideoTable
	model.UserTable
}

type CommentOfVideoInfo struct {
	VideoData
	CommentData `json:"comment"`
}

func (d *Database) GetCommentsWithVideoFromUserInTopic(uid uint, topic string) []CommentGroupByVideo {
	cwd := make([]CommentsWithVideoQuery, 0)
	query := `
	SELECT c.*, v.*, u.*
	FROM comment_tables c
	INNER JOIN video_tables v ON c.video_avid = v.avid
	INNER JOIN user_tables u ON c.owner = u.uid
	WHERE c.owner = ?`
	if topic == "" || topic == "all" || topic == "*" {
		d.db.Raw(query, uid).Scan(&cwd)
	} else {
		query += " AND v.topic = ?"
		d.db.Raw(query, uid, topic).Scan(&cwd)
	}
	return commentsWithVideoQueryToCommentGroupByVideo(cwd)
}

func (d *Database) GetCommentWithVideoByUserFilter(uid uint, name string, rangeType string, ranageData string) []CommentGroupByVideo {
	cwd := make([]CommentsWithVideoQuery, 0)
	query := ""
	switch rangeType {
	case "video":
		query = fmt.Sprintf(" AND v.bvid IN (%s)", ranageData)
	case "topic":
		query = fmt.Sprintf(" AND v.topic IN (%s)", ranageData)
	}
	if uid != 0 && name == "" {
		q := `
		SELECT c.*, v.*, u.*
		FROM comment_tables c
		INNER JOIN video_tables v ON c.video_avid = v.avid
		INNER JOIN user_tables u ON c.owner = u.uid
		WHERE c.owner = ?`
		query = q + query
		d.db.Raw(query, uid).Scan(&cwd)
	} else if uid == 0 && name != "" {
		q := `
		SELECT c.*, v.*, u.*
		FROM user_tables u
		INNER JOIN comment_tables c ON u.uid = c.owner
		INNER JOIN video_tables v ON c.video_avid = v.avid
		WHERE u.username = ?`
		query = q + query
		d.db.Raw(query, name).Scan(&cwd)
	}
	return commentsWithVideoQueryToCommentGroupByVideo(cwd)
}

func commentsWithVideoQueryToCommentGroupByVideo(cwd []CommentsWithVideoQuery) []CommentGroupByVideo {
	userCommentOfVideoInfo := make(map[uint]CommentGroupByVideo)
	for _, item := range cwd {
		current, ok := userCommentOfVideoInfo[item.VideoTable.Avid]
		if !ok {
			current = CommentGroupByVideo{
				VideoData: VideoData{
					Avid:        item.VideoTable.Avid,
					Bvid:        item.VideoTable.Bvid,
					Title:       item.VideoTable.Title,
					Topic:       item.VideoTable.Topic,
					Description: item.VideoTable.Description,
					Cover:       item.VideoTable.Cover,
					Author: Author{
						Uid: item.VideoTable.Owner,
					},
				},
			}
			userCommentOfVideoInfo[item.VideoTable.Avid] = current
		}
		c := CommentData{
			Id:      item.CommentTable.CommentId,
			Content: item.CommentTable.Text,
			Owner: Author{
				Uid:      item.UserTable.Uid,
				Name:     item.UserTable.Username,
				Avatar:   item.UserTable.Avatar,
				Location: item.UserTable.Location,
			},
			Children: nil,
		}
		current.Comments = append(current.Comments, c)
		userCommentOfVideoInfo[item.VideoTable.Avid] = current
	}
	commentsGroup := make([]CommentGroupByVideo, 0)
	for _, v := range userCommentOfVideoInfo {
		commentsGroup = append(commentsGroup, v)
	}
	return commentsGroup
}
