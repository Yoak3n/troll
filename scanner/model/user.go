package model

type UserData struct {
	Uid      uint   `json:"uid"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Avatar   string `json:"avatar"`
}

type UserQuery struct {
	UserTable
	Count int
}
