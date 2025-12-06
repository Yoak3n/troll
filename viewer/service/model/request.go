package model

type VideoTopicRequest struct {
	Topic string `json:"topic"`
	AVID  []uint `json:"avid"`
}
type UpdateTopicRequest struct {
	Topic    string `json:"topic"`
	NewTopic string `json:"new_topic"`
}
