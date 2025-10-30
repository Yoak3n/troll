package model

type CommentData struct {
	Text       string        `json:"text"`
	Author     UserData      `json:"author"`
	Children   []CommentData `json:"children"`
	Rpid       uint          `json:"rpid"`
	Oid        uint          `json:"oid"`
	Like       uint64        `json:"like"`
	NeedExpand bool
}

type CommentResponse struct {
	ResponseCommon `json:"-"`
	Data           CommentResponseData `json:"data"`
}

type LazyCommentResponse struct {
	ResponseCommon `json:"-"`
	Data           LazyCommentData `json:"data"`
}

type LazyCommentData struct {
	Cursor  LazyCommentCursor `json:"cursor"`
	Replies []CommentItem     `json:"replies"`
}

type LazyCommentCursor struct {
	IsBegin         bool            `json:"is_begin"`
	Prev            uint            `json:"prev"`
	Next            uint            `json:"next"`
	IsEnd           bool            `json:"is_end"`
	PaginationReply PaginationReply `json:"pagination_reply"`
}

type PaginationReply struct {
	NextOffset string `json:"next_offset"`
}

type CommentResponseData struct {
	Replies []CommentItem `json:"replies"`
}

type CommentItem struct {
	Rpid         uint                `json:"rpid"`
	Member       CommentMember       `json:"member"`
	Oid          uint                `json:"oid"`
	Mid          uint                `json:"mid"`
	Content      CommentContent      `json:"content"`
	Replies      []CommentItem       `json:"replies"`
	ReplyControl CommentReplyControl `json:"reply_control"`
	Like         uint64              `json:"like"`
}

type CommentContent struct {
	Message  string           `json:"message"`
	Pictures []CommentPicture `json:"pictures"`
	Emote    map[string]any   `json:"emote,omitempty"`
}

type CommentReplyControl struct {
	Location          string `json:"location"`
	SubReplyEntryText string `json:"sub_reply_entry_text"`
}

type CommentPicture struct {
	Src    string  `json:"img_src"`
	Width  int     `json:"img_width"`
	Height int     `json:"img_height"`
	Size   float64 `json:"img_size"`
}
type CommentMember struct {
	Uname  string `json:"uname"`
	Avatar string `json:"avatar"`
}

type SubCommentResponse struct {
	ResponseCommon `json:"-"`
	Data           SubCommentData `json:"data"`
}

type SubCommentData struct {
	Page    SubCommentPage `json:"page"`
	Replies []CommentItem  `json:"replies"`
}

type SubCommentPage struct {
	Count int `json:"count"`
	Num   int `json:"num"`
}
