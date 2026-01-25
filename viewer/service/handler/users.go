package handler

import (
	"net/url"
	"time"

	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/scanner/model"
	viewModel "github.com/Yoak3n/troll/viewer/service/model"
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

func HandlerUserSignedGet(c *fiber.Ctx) error {
	signedUsers, err := controller.GlobalDatabase().GetSignedUserRecord()
	usersList := make([]model.UserData, 0)
	for _, signedUser := range signedUsers {
		usersList = append(usersList, model.UserData{
			Uid:      signedUser.Uid,
			Name:     signedUser.Username,
			Location: signedUser.Location,
			Avatar:   signedUser.Avatar,
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(usersList)
}

func HandlerUserSignPost(c *fiber.Ctx) error {
	req := new(viewModel.UserSignRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	for _, uid := range req.Uids {
		err := controller.GlobalDatabase().CreateSignedRecord(model.SignedUserTable{
			Uid:        uid,
			LastViewed: time.Now(),
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to sign user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user signed",
	})
}

func HandlerUserSignDelete(c *fiber.Ctx) error {
	req := new(viewModel.UserSignRequest)
	if err := c.BodyParser(req); err != nil || len(req.Uids) == 0 {
		uid := c.QueryInt("uid")
		if uid == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body or uid param",
			})
		}
		req.Uids = []uint{uint(uid)}
	}

	for _, uid := range req.Uids {
		err := controller.GlobalDatabase().DeleteSignedUserRecord(uid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to unsign user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user unsigned",
	})
}
