package router

import (
	"geophoto/backend/handler"
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

	// photo
	photo := api.Group("/photo")
	photo.Post("/", handler.CreatePhoto)
	photo.Get("/", handler.GetAllPhotos)
	photo.Get("/:uuid", handler.GetPhoto)
	photo.Put("/:uuid", handler.UpdatePhoto)
	photo.Delete("/:uuid", handler.DeletePhoto)
}
