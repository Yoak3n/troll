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
	Cookie string
}

func Init(dbPath string, dbName string) *Configuration {
	config := &Configuration{
		Auth: Auth{},
	}
	dbConf, err := controller.GlobalDatabase(dbPath, dbName).QueryConfiguration()
	if dbConf == nil || err != nil {
		logger.Logger.Fatal("failed to query configuration, please set config first")
	}
	config.Auth.Cookie = dbConf.Cookie
	config.Proxy = dbConf.Proxy

	Config = config
	return config
}
