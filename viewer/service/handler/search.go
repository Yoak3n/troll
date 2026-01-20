package handler

import (
	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/gofiber/fiber/v2"
)

func HandlerSearhOptionsGet(c *fiber.Ctx) error {
	uid := c.QueryInt("uid", 0)
	name := c.Query("name")
	rangeType := c.Query("rangeType")
	so := controller.SearchOptionResponse{}
	if uid == 0 && name != "" {
		so = controller.GlobalDatabase().GetSearchOptionsByUserName(name, rangeType)
	} else if uid != 0 && name == "" {
		so = controller.GlobalDatabase().GetSearchOptionsByUid(uint(uid), rangeType)
	}
	if len(so.Options) == 0 {
		so = controller.DB.GetSearchOptionsFromAllVideos(rangeType)
	}
	return c.Status(fiber.StatusOK).JSON(so)
}
