package handler

import (
	"github.com/Yoak3n/troll/viewer/config"
	"github.com/gofiber/fiber/v2"
)

func HandlerSettingGet(c *fiber.Ctx) error {
	conf := config.GetConfiguration()
	return c.Status(fiber.StatusOK).JSON(conf)
}

func HandlerSettingUpdate(c *fiber.Ctx) error {
	conf := &config.Configuration{}
	err := c.BodyParser(conf)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	config.UpdateAllConfiguration(conf)
	return c.Status(fiber.StatusOK).JSON(nil)
}
