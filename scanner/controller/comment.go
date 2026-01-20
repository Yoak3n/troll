package controller

import (
	"sync"

	"github.com/Yoak3n/troll/scanner/database"
	"github.com/Yoak3n/troll/scanner/model"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

var DB *Database
var once sync.Once

func NewDatabase(dir string, name string) *Database {
	return &Database{
		db: database.InitDatabase(dir, name),
	}
}

func GlobalDatabase(args ...string) *Database {
	once.Do(func() {
		if len(args) == 2 {
			DB = NewDatabase(args[0], args[1])
		}
	})
	return DB
}

func (d *Database) AddCommentRecord(comments []model.CommentTable) error {
	if len(comments) == 0 {
		return nil
	}
	return d.db.Save(&comments).Error
}

func (d *Database) QueryCommentRecord(comment model.CommentTable) (*model.CommentTable, error) {
	return &comment, d.db.First(&comment).Error
}

func (d *Database) QueryCommentList(video model.VideoTable) ([]model.CommentTable, error) {
	var ret []model.CommentTable
	err := d.db.Where("video_avid  = ?", video.Avid).Find(&ret).Error
	return ret, err
}

func (d *Database) QueryUserCommentList(user model.UserTable) ([]model.CommentTable, error) {
	var ret []model.CommentTable
	err := d.db.Where("owner = ?", user.Uid).Find(&ret).Error
	return ret, err
}

func (d *Database) AddUserRecord(users []model.UserTable) error {
	if len(users) == 0 {
		return nil
	}
	return d.db.Save(&users).Error
}

func (d *Database) QueryUserCommentsListInTopic(topic string, username string) ([]model.CommentTable, error) {
	comments := make([]model.CommentTable, 0)
	query := `
		SELECT c.*
		FROM comment_tables c
		INNER JOIN video_tables v ON c.video_avid = v.avid
		INNER JOIN user_tables u ON c.owner = u.uid
		WHERE u.username = ? AND v.topic = ?
    `

	err := d.db.Raw(query, username, topic).Scan(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (d *Database) QuerySimilarComments(topic string, n int) ([]model.SimilarCommentResult, error) {
	comments := make([]model.SimilarCommentResult, 0)

	// 使用 Where 和 Joins 方法构建查询，而不是直接使用 Raw SQL
	subQuery := d.db.Table("comment_tables AS c").
		Select("SUBSTR(c.text, 1, 30) AS content_prefix, COUNT(*) AS similar_count, GROUP_CONCAT(c.comment_id, ', ') AS comment_ids").
		Joins("INNER JOIN video_tables v ON c.video_avid = v.avid").
		Where("c.text IS NOT NULL AND LENGTH(c.text) >= 10 AND v.topic = ?", topic).
		Group("content_prefix").
		Having("COUNT(*) > 1")

	query := d.db.Table("(?) AS similar_groups", subQuery).
		Order("similar_count DESC, content_prefix").
		Limit(n)

	err := query.Scan(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

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
	WHERE c.video_avid = ?  AND c.parent_comment = 0 AND c.deleted_at IS NULL`
	rootQuery := make([]CommentQuery, 0)
	if err := d.db.Raw(query1, avid).Scan(&rootQuery).Error; err != nil {
		return CommentGroupByVideo{}
	}
	query2 := `
	SELECT c.*, u.*
	FROM comment_tables c
	INNER JOIN user_tables u ON c.owner = u.uid
	WHERE c.video_avid = ?  AND c.parent_comment != 0 AND c.deleted_at IS NULL`
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

func (d *Database) SearchCommentWithKeyword(keyword string) []CommentData {
	query := `
	SELECT c.*,u.*
	FROM comment_tables c
	INNER JOIN user_tables u ON c.owner = u.uid
	WHERE c.text LIKE ? AND c.deleted_at IS NULL`
	result := make([]CommentQuery, 0)
	if err := d.db.Raw(query, "%"+keyword+"%").Scan(&result).Error; err != nil {
		return nil
	}
	comments := make([]CommentData, 0, len(result))
	for _, item := range result {
		comments = append(comments, CommentData{
			Id:      item.CommentId,
			Content: item.Text,
			Owner: Author{
				Uid:      item.Uid,
				Name:     item.Username,
				Avatar:   item.Avatar,
				Location: item.Location,
			},
		})
	}
	return comments
}
