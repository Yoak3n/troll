package controller

import "github.com/Yoak3n/troll/scanner/model"

type TopicsData struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (d *Database) GetAllTopicsList() []TopicsData {
	result := make([]model.VideoTable, 0)
	d.db.Find(&result)
	topics := make(map[string]int64)
	for _, video := range result {
		topics[video.Topic]++
	}
	ret := make([]TopicsData, 0)
	for k, v := range topics {
		ret = append(ret, TopicsData{
			Name:  k,
			Count: v,
		})
	}
	return ret
}
