package model

import (
	"time"

	"gorm.io/gorm"
)

type ConfigurationTable struct {
	Cookie string
	Proxy  string
	gorm.Model
}

type UserTable struct {
	Uid       uint   `json:"uid" gorm:"primaryKey"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type VideoTable struct {
	Avid        uint   `json:"avid" gorm:"primaryKey"`
	Bvid        string `json:"bvid"`
	Title       string `json:"title"`
	Topic       string
	Description string `json:"description"`
	Owner       uint   `json:"owner"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type CommentTable struct {
	Text          string `json:"text"`
	Owner         uint
	VideoAvid     uint
	ParentComment uint
	CommentId     uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type SimilarCommentResult struct {
	Text       string
	Count      int
	CommentIds string
	Example1   string
	Example2   string
}
