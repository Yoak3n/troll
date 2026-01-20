package app

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/Yoak3n/troll/scanner/controller"
	"github.com/Yoak3n/troll/viewer/consts"
	"github.com/Yoak3n/troll/viewer/service/handler"
	"github.com/Yoak3n/troll/viewer/service/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	handler2 "github.com/Yoak3n/troll/scanner/package/handler"
)

func initServices() {
	controller.GlobalDatabase(consts.TrollPath, "troll")
	ws.InitWebsocketHub()
	handler.InitHandlerState()
	handler2.Init(consts.TrollPath, "troll")
}

func initRouter(files fs.FS, port int) error {
	app := fiber.New(fiber.Config{
		AppName: "Troll Viewer",
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))
	setupRoutes(app)

	subFS, err := fs.Sub(files, "dist")
	if err != nil {
		return err
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   http.FS(subFS),
		Browse: false,
		Index:  "index.html",
		MaxAge: 3600,
	}))

	app.Listen(fmt.Sprintf(":%d", port))
	return nil
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return c.Status(403).SendString("Request origin not allowed")
	})
	v1.Get("/ws/:id", websocket.New(func(conn *websocket.Conn) {
		id := conn.Params("id")
		if id != "" {
			ws.Hub.Register(id, conn)
		} else {
			conn.WriteMessage(400, []byte("need a id"))
		}
	}))
	v1.Post("/task/video/refresh", handler.HandlerVideoRefreshPost)
	setupTopicsRoutes(v1)
	setupVideosRoutes(v1)
	setupUserRoutes(v1)
	setupSearchRoutes(v1)
	setupSettingRoutes(v1)
	setupCommentRoutes(v1)
	setupStatisticsRoutes(v1)
}

func setupTopicsRoutes(group fiber.Router) {
	topics := group.Group("/topics")
	topics.Get("/list", handler.HandlerTopicsGet).Name("list")
	topics.Get("/:topicName/videos", handler.HandlerTopicVideosGet).Name("videos")
	topics.Post("/update", handler.HandlerTopicUpdate).Name("update")
	topics.Delete("/:topicName", handler.HandlerTopicDelete).Name("delete")
}

func setupVideosRoutes(group fiber.Router) {
	videos := group.Group("/videos")
	videos.Get("/:avid/comments", handler.HandlerVideoCommentsGet).Name("comments")
	videos.Post("/topic", handler.HandlerVideoTopicPost).Name("topic")
	videos.Delete("/", handler.HandlerVideosDelete).Name("delete")
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

func setupSettingRoutes(group fiber.Router) {
	setting := group.Group("/setting").Name("setting.")
	setting.Get("/", handler.HandlerSettingGet).Name("Get")
	setting.Post("/", handler.HandlerSettingUpdate).Name("Update")
}

func setupCommentRoutes(group fiber.Router) {
	comments := group.Group("/comments").Name("comments.")
	comments.Get("/search", handler.HandleCommentSearchWithKeyword).Name("search")
}

func setupStatisticsRoutes(group fiber.Router) {
	stats := group.Group("/statistics").Name("statistics.")
	stats.Get("/", handler.HandlerDashboardStatsGet).Name("dashboard")
}
