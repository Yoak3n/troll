package handler

import (
	"net/url"

	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/gofiber/fiber/v2"
)

func HandleCommentSearchWithKeyword(c *fiber.Ctx) error {
	keyword, _ := url.QueryUnescape(c.Query("keyword"))
	output := controller.GlobalDatabase().SearchCommentWithKeyword(keyword)
	if output == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no comment with this keyword",
		})
	}
	return c.Status(fiber.StatusOK).JSON(output)
}
