package response

import (
	"github.com/gofiber/fiber/v2"
)

func InternalServerError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Internal server error"})
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Not found"})
}

func InvalidInput(c *fiber.Ctx, reason interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": reason})
}

func Data(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": data})
}

func Message(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": msg})
}

func RecordCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Record created", "data": data})
}

func RecordUpdated(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Record updated", "data": data})
}

func RecordDeleted(c *fiber.Ctx) error {
	return Message(c, "Record deleted")
}
