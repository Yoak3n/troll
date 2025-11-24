package handler

import (
	"strings"

	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/troll/scanner/internal/config"
)

func Init(dir string, name string) {
	logger.Init()
	config.Init(dir, name)
	InitAccountLimiter()
}

type Handler struct {
	title string
	topic string
	cache string
	bvid  string
	avid  int64
}

func NewHandler(cache string, title string, topic string, bvid string, avid int64) *Handler {
	hub := &Handler{
		title: title,
		topic: topic,
		cache: cache,
		bvid:  bvid,
		avid:  avid,
	}
	return hub
}

func (h *Handler) Run() {
	if h.topic != "" {
		NewTopic(h.cache, h.title, strings.Split(h.topic, ","))
		return
	}
	if h.bvid != "" || h.avid != -1 {
		NewVideo(h.cache, h.title, h.bvid, h.avid)
		return
	}

}
