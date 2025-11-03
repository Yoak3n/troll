package model

type VideoData struct {
	Avid        uint          `json:"avid"`
	Bvid        string        `json:"bvid"`
	Title       string        `json:"title"`
	Cover       string        `json:"cover"`
	Description string        `json:"description"`
	Owner       UserData      `json:"owner"`
	Comments    []CommentData `json:"comments"`
}

type VideoInfoResponse struct {
	ResponseCommon `json:"-"`
	Data           VideoInfoData `json:"data"`
}

type VideoInfoData struct {
	Aid         uint           `json:"aid"`
	Bvid        string         `json:"bvid"`
	Title       string         `json:"title"`
	Pic         string         `json:"pic"`
	Description string         `json:"desc"`
	Owner       VideoInfoOwner `json:"owner"`
}

type VideoInfoOwner struct {
	Mid  uint   `json:"mid"`
	Name string `json:"name"`
}
