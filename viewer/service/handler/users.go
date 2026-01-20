package handler

import (
	"net/url"

	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/gofiber/fiber/v2"
)

func HandlerUserCommentGet(c *fiber.Ctx) error {
	uid, err := c.ParamsInt("uid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uid parameter",
		})
	}
	topic := c.Query("topicName")
	topic, _ = url.QueryUnescape(topic)
	comments := controller.GlobalDatabase().GetCommentsWithVideoFromUserInTopic(uint(uid), topic)
	return c.Status(fiber.StatusOK).JSON(comments)
}

func HandlerUserCommentsFilter(c *fiber.Ctx) error {
	uid := c.QueryInt("uid")
	name, _ := url.QueryUnescape(c.Query("name"))
	rangeType, _ := url.QueryUnescape(c.Query("rangeType"))
	rangeData, _ := url.QueryUnescape(c.Query("rangeData"))
	output := controller.GlobalDatabase().GetCommentWithVideoByUserFilter(uint(uid), name, rangeType, rangeData)
	return c.Status(fiber.StatusOK).JSON(output)
}
