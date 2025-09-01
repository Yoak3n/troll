package handler

import (
	"strings"

	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/troll-scanner/internal/config"
)

func init() {
	logger.Init()
	config.Init()
}

type Handler struct {
	Title string
	Topic string
	cache string
}

func NewHandler(cache string, title string, topic string) *Handler {
	hub := &Handler{
		Title: title,
		Topic: topic,
		cache: cache,
	}
	return hub
}

func (h *Handler) Run() {
	NewTopic(h.cache, h.Title, strings.Split(h.Topic, ","))
}
