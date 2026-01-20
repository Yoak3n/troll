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

type DashboardStats struct {
	Topics   int64 `json:"topics"`
	Videos   int64 `json:"videos"`
	Users    int64 `json:"users"`
	Comments int64 `json:"comments"`
}

func (d *Database) GetDashboardStats() DashboardStats {
	var topics int64
	var videos int64
	var users int64
	var comments int64

	d.db.Raw("SELECT COUNT(DISTINCT topic) FROM video_tables").Scan(&topics)
	d.db.Model(&model.VideoTable{}).Count(&videos)
	d.db.Model(&model.UserTable{}).Count(&users)
	d.db.Model(&model.CommentTable{}).Count(&comments)

	return DashboardStats{
		Topics:   topics,
		Videos:   videos,
		Users:    users,
		Comments: comments,
	}
}

func (d *Database) UpdateTopic(topic string, newTopic string) error {
	d.db.Model(&model.VideoTable{}).Where("topic = ?", topic).Update("topic", newTopic)
	return nil
}

func (d *Database) DeleteTopic(topic string) error {
	avids := make([]uint, 0)
	if err := d.db.Model(&model.VideoTable{}).Where("topic = ?", topic).Pluck("avid", &avids).Error; err != nil {
		return err
	}
	if len(avids) > 0 {
		if err := d.db.Where("video_avid IN ?", avids).Delete(&model.CommentTable{}).Error; err != nil {
			return err
		}
	}
	if err := d.db.Delete(&model.VideoTable{}, "topic = ?", topic).Error; err != nil {
		return err
	}
	return nil
}
