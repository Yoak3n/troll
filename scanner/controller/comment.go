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

func (d *Database) AddVideoRecord(video model.VideoTable) error {
	return d.db.Save(&video).Error
}

func (d *Database) QueryVideoRecord(video model.VideoTable) (*model.VideoTable, error) {
	return &video, d.db.First(&video).Error
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

func (d *Database) QueryConfigurationCookie() (*model.ConfigurationTable, error) {
	c := &model.ConfigurationTable{
		Type: "cookie",
	}
	err := d.db.First(c).Error
	return c, err
}

func (d *Database) QueryConfigurationProxy() (*model.ConfigurationTable, error) {
	c := &model.ConfigurationTable{
		Type: "proxy",
	}
	err := d.db.First(c).Error
	return c, err
}

func (d *Database) QueryConfiguration() ([]model.ConfigurationTable, error) {
	confs := make([]model.ConfigurationTable, 0)
	err := d.db.Find(&confs).Error
	return confs, err
}

func (d *Database) UpdateConfiguration(c *model.ConfigurationTable) error {
	return d.db.Save(c).Error
}
