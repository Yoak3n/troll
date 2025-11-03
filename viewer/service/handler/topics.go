package handler

import (
	"encoding/json"
	"net/url"

	"github.com/Yoak3n/troll/viewer/service/controller"
	"github.com/gofiber/fiber/v2"
)

func HandlerTopicsGet(c *fiber.Ctx) error {
	db := controller.GlobalDatabase()
	topics := db.GetAllTopicsList()
	buf, err := json.Marshal(topics)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return c.Send(buf)
}

func HandlerTopicVideosGet(c *fiber.Ctx) error {
	topicName := c.Params("topicName")
	de, err := url.QueryUnescape(topicName)
	if err == nil {
		topicName = de
	}
	db := controller.GlobalDatabase()
	videos := db.GetVideosByTopic(de)
	buf, err := json.Marshal(videos)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return c.Send(buf)
}
