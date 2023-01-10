package response

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func PlaceHolder(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true, "message": "in development", "data": nil})
}

func InternalServerError(c *fiber.Ctx, err *error) error {
	if err != nil {
		// log.Println(string(debug.Stack()))
		log.Println(*err)
	}
	return c.Status(500).JSON(fiber.Map{"success": false, "message": "Internal server error"})
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data not found", "data": nil})
}

func Data(c *fiber.Ctx, data interface{}) error {
	return c.Status(200).JSON(fiber.Map{"success": true, "data": data})
}

func Message(c *fiber.Ctx, msg string) error {
	return c.Status(200).JSON(fiber.Map{"success": true, "message": msg})
}

func DataMessage(c *fiber.Ctx, data interface{}, msg string) error {
	return c.Status(200).JSON(fiber.Map{"success": true, "data": data, "message": msg})
}
