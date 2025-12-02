package model

type SearchResponse struct {
	ResponseCommon `json:"-"`
	Data           SearchResponseData `json:"data"`
}
type SearchResponseData struct {
	Result []SearchItem `json:"result"`
}

type SearchItem struct {
	Typ         string `json:"type"`
	Id          uint   `json:"id"`
	Author      string `json:"author"`
	Pic         string `json:"pic"`
	Mid         uint   `json:"mid"`
	Title       string `json:"title"`
	Aid         uint   `json:"aid"`
	Bvid        string `json:"bvid"`
	Description string `json:"description"`
	Review      int    `json:"review"`
}
