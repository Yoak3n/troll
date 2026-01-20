package config

import (
	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/model"
	"gorm.io/gorm"
)

type Configuration struct {
	Proxies  []ConfigurationItem `json:"proxies"`
	Cookies  []ConfigurationItem `json:"cookies"`
	Interval RequestInterval     `json:"interval"`
}

type ConfigurationItem struct {
	Id   uint   `json:"id"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type RequestInterval struct {
	Basic  int `json:"basic"`
	Random int `json:"random"`
}

var Conf *Configuration

func GetConfiguration() *Configuration {
	confs := controller.DB.QueryConfigurations()
	configuration := &Configuration{
		Proxies: make([]ConfigurationItem, 0),
		Cookies: make([]ConfigurationItem, 0),
		// 暂不支持修改
		Interval: RequestInterval{
			Basic:  1,
			Random: 3,
		},
	}
	for _, conf := range confs {
		item := ConfigurationItem{
			Id:   conf.ID,
			Data: conf.Data,
			Type: conf.Type,
		}
		switch conf.Type {
		case "cookie":
			configuration.Cookies = append(configuration.Cookies, item)
		case "proxy":
			configuration.Proxies = append(configuration.Proxies, item)
		}
	}
	Conf = configuration
	return configuration
}

func getConfigurationMap() map[uint]string {
	currentConfigurationMap := make(map[uint]string)
	for _, cookies := range Conf.Cookies {
		currentConfigurationMap[cookies.Id] = cookies.Data
	}
	for _, proxy := range Conf.Proxies {
		currentConfigurationMap[proxy.Id] = proxy.Data
	}
	return currentConfigurationMap
}

func UpdateAllConfiguration(conf *Configuration) {
	currentConfigurationMap := getConfigurationMap()
	newMap := make(map[uint]string)
	diffConfiguration := make([]model.ConfigurationTable, 0)
	for _, cookie := range conf.Cookies {
		diffConfiguration = append(diffConfiguration, model.ConfigurationTable{
			Model: gorm.Model{
				ID: cookie.Id,
			},
			Type: cookie.Type,
			Data: cookie.Data,
		})
		newMap[cookie.Id] = cookie.Data
	}
	for _, proxy := range conf.Proxies {
		newMap[proxy.Id] = proxy.Data
		diffConfiguration = append(diffConfiguration, model.ConfigurationTable{
			Model: gorm.Model{
				ID: proxy.Id,
			},
			Data: proxy.Data,
			Type: proxy.Type,
		})

	}
	controller.DB.UpdateConfiguration(diffConfiguration)

	deleteSlice := make([]uint, 0)
	for k := range currentConfigurationMap {
		if _, ok := newMap[k]; !ok {
			deleteSlice = append(deleteSlice, k)
		}
	}
	controller.DB.DeleteConfiguration(deleteSlice)
}
