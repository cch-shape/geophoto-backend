package handler

import (
	"geophoto/backend/utils/response"
	"github.com/gofiber/fiber/v2"
	"os"
)

func Hello(c *fiber.Ctx) error {
	return response.Message(c, "Hello")
}

func GetAllConfigs(c *fiber.Ctx) error {
	return response.Data(c, fiber.Map{
		"image_path": os.Getenv("IMAGE_PATH"),
	})
}
