package middleware

import (
	"database/sql"
	"errors"
	"geophoto/backend/utils/response"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if err == sql.ErrNoRows {
		return response.NotFound(c)
	} else if err.Error() == "validation failed" {
		return response.InvalidInput(c, c.Locals("reason"))
	} else if err.Error() == "there is no uploaded file associated with the given key" {
		return response.InvalidInput(c, "file is required")
	} else {
		var e *fiber.Error
		if errors.As(err, &e) {
			return c.Status(e.Code).SendString(err.Error())
		}
	}

	return response.InternalServerError(c)
}
