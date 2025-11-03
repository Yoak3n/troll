package handler

import (
	"encoding/json"
	"fmt"

	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/internal/util"
	"github.com/Yoak3n/troll/scanner/model"
)

const UserApi = "https://api.bilibili.com/x/space/wbi/acc/info"

func AddUserByUid(uid uint) {
	logger.Logger.Printf("====Fetch user:%d begining====", uid)
	resp := util.RequestGetWithAll(UserApi + fmt.Sprintf("?mid=%d", uid))
	type UserResponse struct {
		Code int `json:"code"`
		Data struct {
			Mid  uint   `json:"mid"`
			Name string `json:"name"`
			Face string `json:"face"`
		} `json:"data"`
	}
	var userResponse UserResponse
	err := json.Unmarshal(resp, &userResponse)
	if err != nil {
		logger.Logger.Errorln(err)
		return
	}
	var user model.UserTable
	user.Uid = userResponse.Data.Mid
	user.Username = userResponse.Data.Name
	user.Avatar = userResponse.Data.Face

	err = controller.GlobalDatabase().AddUserRecord([]model.UserTable{user})
	if err != nil {
		logger.Logger.Errorln(err)
		return
	}
}
