package handler

import (
	"encoding/json"
	"net/url"

	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/viewer/service/model"
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

func HandlerDashboardStatsGet(c *fiber.Ctx) error {
	stats := controller.GlobalDatabase().GetDashboardStats()
	buf, err := json.Marshal(stats)
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

func HandlerTopicUpdate(c *fiber.Ctx) error {
	var req model.UpdateTopicRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}
	db := controller.GlobalDatabase()
	err = db.UpdateTopic(req.Topic, req.NewTopic)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return c.SendString("OK")
}

func HandlerTopicDelete(c *fiber.Ctx) error {
	topicName := c.Params("topicName")
	de, err := url.QueryUnescape(topicName)
	if err == nil {
		topicName = de
	}
	db := controller.GlobalDatabase()
	err = db.DeleteTopic(topicName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return c.SendString("OK")
}
