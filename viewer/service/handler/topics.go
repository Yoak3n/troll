package handler

import (
	"encoding/json"

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
