package handler

import (
	"fmt"
	"geophoto/backend/database"
	"geophoto/backend/model"
	"geophoto/backend/utils"
	"github.com/gofiber/fiber/v2"
)

func GetAllPhotosWithJWT(c *fiber.Ctx) error {
	var photos []model.Photo
	err := database.DB.Select(&photos, utils.SelectStmt(&photos, database.TableName["Photo"]))
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Internal server error"})
	}

	return c.JSON(fiber.Map{"success": true, "data": photos})
}
