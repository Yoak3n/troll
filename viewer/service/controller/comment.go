package controller

import "github.com/Yoak3n/troll/scanner/model"

type CommentData struct {
	Id       uint          `json:"id"`
	Content  string        `json:"content"`
	Owner    Author        `json:"owner"`
	Children []CommentData `json:"children,omitempty"`
}

type CommentQuery struct {
	model.CommentTable
	model.UserTable
}

type CommentGroupByVideo struct {
	VideoData
	Comments []CommentData `json:"comments"`
}

func (d *Database) GetCommentsByVideo(avid uint) CommentGroupByVideo {

	videoData := d.GetVideoByAvid(avid)
	query1 := `
	SELECT c.*,u.*
	FROM comment_tables c
	INNER JOIN user_tables u ON c.owner = u.uid
	WHERE c.video_avid = ?  AND c.parent_comment = 0`
	rootQuery := make([]CommentQuery, 0)
	if err := d.db.Raw(query1, avid).Scan(&rootQuery).Error; err != nil {
		return CommentGroupByVideo{}
	}
	query2 := `
	SELECT c.*, u.*
	FROM comment_tables c
	INNER JOIN user_tables u ON c.owner = u.uid
	WHERE c.video_avid = ?  AND c.parent_comment != 0`
	subQuery := make([]CommentQuery, 0)
	if err := d.db.Raw(query2, avid).Scan(&subQuery).Error; err != nil {
		return CommentGroupByVideo{}
	}
	// rootComments := make([]model.CommentTable, 0)
	// subComments := make([]model.CommentTable, 0)
	commentMap := make(map[uint]CommentData)
	for _, item := range rootQuery {
		commentMap[item.CommentId] = CommentData{
			Id:      item.CommentId,
			Content: item.Text,
			Owner: Author{
				Uid:      item.Uid,
				Name:     item.Username,
				Avatar:   item.Avatar,
				Location: item.Location,
			},
			Children: make([]CommentData, 0),
		}
	}
	for _, item := range subQuery {
		parentComment, ok := commentMap[item.ParentComment]
		if ok {
			parentComment.Children = append(parentComment.Children, CommentData{
				Id:      item.CommentId,
				Content: item.Text,
				Owner: Author{
					Uid:      item.Uid,
					Name:     item.Username,
					Avatar:   item.Avatar,
					Location: item.Location,
				},
			})
			commentMap[item.ParentComment] = parentComment
		}
	}
	commentsGroupByVideo := CommentGroupByVideo{
		VideoData: videoData,
		Comments:  make([]CommentData, 0, len(commentMap)),
	}
	for _, comment := range commentMap {
		commentsGroupByVideo.Comments = append(commentsGroupByVideo.Comments, comment)
	}
	return commentsGroupByVideo
}
