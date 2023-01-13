package router

import (
	"geophoto/backend/handler"
	"geophoto/backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"path/filepath"
)

func SetupRoutes(app *fiber.App) {
	// static
	app.Static(
		filepath.Join("/", os.Getenv("IMAGE_PATH")),
		filepath.Join(".", os.Getenv("IMAGE_PATH")),
	)

	// api
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// auth
	api.Post("/login", handler.Login)
	api.Post("/ask/verification", handler.AskVerificationCode)

	// photo
	photo := api.Group("/photo")
	photo.Post("/", middleware.Protected(), handler.CreatePhoto)
	photo.Get("/", middleware.Protected(), handler.GetAllPhotos)
	photo.Get("/:uuid", middleware.Protected(), handler.GetPhoto)
	photo.Put("/:uuid", middleware.Protected(), handler.UpdatePhoto)
	photo.Delete("/:uuid", middleware.Protected(), handler.DeletePhoto)

	// user
	user := api.Group("/user")
	user.Get("/", middleware.Protected(), handler.GetUserSelf)
	user.Put("/", middleware.Protected(), handler.UpdateUserSelf)
}
