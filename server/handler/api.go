package handler

import (
	"geophoto/backend/utils/response"
	"github.com/gofiber/fiber/v2"
)

func Hello(c *fiber.Ctx) error {
	return response.Message(c, "Hello")
}
