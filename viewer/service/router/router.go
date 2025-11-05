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
	setupVideosRoutes(v1)
	setupUserRoutes(v1)
	setupSearchRoutes(v1)
}

func setupTopicsRoutes(group fiber.Router) {
	topics := group.Group("/topics")
	topics.Get("/list", handler.HandlerTopicsGet).Name("list")
	topics.Get("/:topicName/videos", handler.HandlerTopicVideosGet).Name("videos")
}

func setupVideosRoutes(group fiber.Router) {
	videos := group.Group("/videos")
	videos.Get("/:avid/comments", handler.HandlerVideoCommentsGet).Name("comments")
}

func setupUserRoutes(group fiber.Router) {
	users := group.Group("/users").Name("users.")
	users.Get("/:uid/comments", handler.HandlerUserCommentGet).Name("user.comments")
	users.Get("/filter/coments", handler.HandlerUserCommentsFilter).Name("comments")
}

func setupSearchRoutes(group fiber.Router) {
	search := group.Group("/search").Name("search.")
	search.Get("/options", handler.HandlerSearhOptionsGet).Name("options")
}
