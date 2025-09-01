package model

type VideoData struct {
	Avid        uint          `json:"avid"`
	Bvid        string        `json:"bvid"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Owner       UserData      `json:"owner"`
	Comments    []CommentData `json:"comments"`
}
