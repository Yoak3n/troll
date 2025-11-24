package config

import (
	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/troll/scanner/controller"
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
			config.Auth.Cookie = append(config.Auth.Cookie, conf.Data)
		}
		if conf.Type == "proxy" {
			config.Proxy = conf.Data
		}
	}

	Config = config
	return config
}
