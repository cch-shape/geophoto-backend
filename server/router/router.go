package router

import (
	"geophoto/backend/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)
	api.Get("/configs", handler.GetAllConfigs)

	// Photo
	photo := api.Group("/photo")
	photo.Post("/", handler.CreatePhoto)
	photo.Get("/", handler.GetAllPhotos)
	photo.Get("/:uuid", handler.GetPhoto)
	photo.Patch("/:uuid", handler.UpdatePhoto)
	photo.Delete("/:uuid", handler.DeletePhoto)
}
