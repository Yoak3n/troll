package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/internal/util"
	"github.com/Yoak3n/troll/scanner/model"
	util2 "github.com/Yoak3n/troll/scanner/package/util"
)

const UserApi = "https://api.bilibili.com/x/space/wbi/acc/info"

func AddUserByUid(uid uint) {
	url := util2.AppendParamsToUrl(UserApi, map[string]string{
		"mid": fmt.Sprintf("%d", uid),
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	accountID, cookie := accountLimiter.GetAccount(ctx)
	resp := util.RequestGetWithAll(url, cookie)
	cancel()
	type UserResponse struct {
		Code int `json:"code"`
		Data struct {
			Mid  uint   `json:"mid"`
			Name string `json:"name"`
			Face string `json:"face"`
		} `json:"data"`
	}
	userResponse := &UserResponse{}
	err := json.Unmarshal(resp, userResponse)
	if err != nil {
		accountLimiter.Penalize(accountID)
		logger.Logger.Errorln(err)
		return
	}
	var user model.UserTable
	user.Uid = userResponse.Data.Mid
	user.Username = userResponse.Data.Name
	user.Avatar = userResponse.Data.Face

	err = controller.GlobalDatabase().AddUserRecord([]model.UserTable{user})
	if err != nil {
		accountLimiter.Penalize(accountID)
		logger.Logger.Errorln(err)
		return
	}
	accountLimiter.Reward(accountID)
}
