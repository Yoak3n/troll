package controller

import "github.com/Yoak3n/troll/scanner/model"

func (d *Database) AddVideoRecord(video model.VideoTable) error {
	return d.db.Save(&video).Error
}

func (d *Database) QueryVideoRecord(video model.VideoTable) (*model.VideoTable, error) {
	return &video, d.db.First(&video).Error
}

func (d *Database) UpdateTopicofVideos(ids []uint, topic string) error {
	return d.db.Model(&model.VideoTable{}).Where("avid IN ?", ids).Update("topic", topic).Error
}
