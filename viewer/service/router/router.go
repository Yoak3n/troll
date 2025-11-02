package router

import (
	"github.com/Yoak3n/troll/viewer/consts"
	"github.com/Yoak3n/troll/viewer/service/controller"
	"github.com/Yoak3n/troll/viewer/service/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouter() {
	controller.GlobalDatabase(consts.TrollPath, "troll")
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	setupRoutes(app)
	app.Listen(":10420")
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	setupTopicsRoutes(v1)
}

func setupTopicsRoutes(group fiber.Router) {
	topics := group.Group("/topics")
	topics.Get("/", handler.HandlerTopicsGet).Name("get")
}
