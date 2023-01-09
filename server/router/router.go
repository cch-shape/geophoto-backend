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

	// Photo
	photo := api.Group("/photo")
	photo.Get("/", handler.GetAllPhotosWithJWT)
}
