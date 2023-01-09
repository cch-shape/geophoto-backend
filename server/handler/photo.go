package handler

import (
	"fmt"
	"geophoto/backend/database"
	"geophoto/backend/model"
	"github.com/gofiber/fiber/v2"
)

func GetAllPhotosWithJWT(c *fiber.Ctx) error {
	var photos []model.Photo
	stmt := fmt.Sprintf("SELECT id, user_id, photo_url, description, timestamp, X(coordinates), Y(coordinates) FROM %s", database.TableName["Photo"])
	err := database.DB.Select(&photos, stmt)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Internal server error"})
	}

	return c.JSON(fiber.Map{"success": true, "data": photos})
}
