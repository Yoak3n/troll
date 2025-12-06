package handler

import (
	"sort"

	"github.com/Yoak3n/troll/viewer/service/controller"
	"github.com/Yoak3n/troll/viewer/service/model"
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
	sortType := c.Query("sort", "children_desc")
	if sortType == "children_desc" && len(comments.Comments) > 1 {
		sort.Slice(comments.Comments, func(i, j int) bool {
			il := len(comments.Comments[i].Children)
			jl := len(comments.Comments[j].Children)
			return il > jl
		})
	}
	return c.Status(fiber.StatusOK).JSON(comments)

}

func HandlerVideoTopicPost(c *fiber.Ctx) error {
	req := &model.VideoTopicRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	topic := req.Topic
	if topic == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "topic is empty",
		})
	}
	err = controller.GlobalDatabase().UpdateTopicOfVideos(req.AVID, topic)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update topic",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "topic updated",
	})
}
