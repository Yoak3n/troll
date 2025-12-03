package config

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/model"
	"gorm.io/gorm"
)

type Configuration struct {
	Proxy string
	Auth  Auth
}

var Config *Configuration

type Auth struct {
	Cookie []string
}

func Init(dbPath string, dbName string) *Configuration {
	config := &Configuration{
		Auth: Auth{},
	}
	dbConfs, err := controller.GlobalDatabase(dbPath, dbName).QueryConfiguration()
	if dbConfs == nil || err != nil {
		logger.Logger.Fatal("failed to query configuration, please set config first")
	}
	for _, conf := range dbConfs {
		if conf.Type == "cookie" {
			if !CheckCookie(conf.Data) {
				logger.Logger.Fatalf("cookie %s is invalid", conf.Data)
				controller.GlobalDatabase().UpdateConfiguration(&model.ConfigurationTable{
					Model: gorm.Model{
						ID: conf.ID,
					},
					Invalid: true,
				})
				continue
			}
			config.Auth.Cookie = append(config.Auth.Cookie, conf.Data)
		}
		if conf.Type == "proxy" {
			config.Proxy = conf.Data
		}
	}

	Config = config
	if len(config.Auth.Cookie) == 0 {
		logger.Logger.Fatal("no valid cookie found")
	}
	return config
}

const CookieCheckUri = "https://passport.bilibili.com/x/passport-login/web/cookie/info"

type CookieInfoResponse struct {
	Code int `json:"code"`
}

func CheckCookie(cookie string) bool {
	req, _ := http.NewRequest("GET", CookieCheckUri, nil)
	req.Header.Set("Cookie", cookie)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorf("CheckCookie err: %v", err)
		return false
	}
	defer res.Body.Close()
	resBuf, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Logger.Errorf("CheckCookie err: %v", err)
		return false
	}
	response := &CookieInfoResponse{}
	err = json.Unmarshal(resBuf, response)
	if err != nil || response.Code != 0 {
		logger.Logger.Errorf("CheckCookie err: %v %d", err, response.Code)
		return false
	}
	return true
}
