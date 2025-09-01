package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Proxy string
	Auth  Auth
}

var Config *Configuration

type Auth struct {
	Cookie string
}

func Init() *Configuration {
	config := &Configuration{
		Auth: Auth{},
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	config.Auth.Cookie = viper.GetString("auth.cookie")
	config.Proxy = viper.GetString("proxy")

	Config = config
	return config
}
