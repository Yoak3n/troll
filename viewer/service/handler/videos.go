package handler

import (
	"github.com/Yoak3n/troll/viewer/service/controller"
	"github.com/gofiber/fiber/v2"
)

func HandlerVideoCommentsGet(c *fiber.Ctx) error {
	avid, err := c.ParamsInt("avid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid avid parameter",
		})
	}
	comments := controller.GlobalDatabase().GetCommentsByVideo(uint(avid))
	return c.Status(fiber.StatusOK).JSON(comments)

}
