package dto

import (
	"github.com/Yoak3n/troll/scanner/model"
)

type VideoDataOutput struct {
	VideoID          string `json:"video_id"`
	Count            uint   `json:"count"`
	*model.VideoData `json:"data,omitempty"`
}
