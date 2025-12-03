package handler

import (
    "sort"

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
