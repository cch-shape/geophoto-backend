package handler

import (
	"geophoto/backend/model"
	"geophoto/backend/utils/response"
	"github.com/gofiber/fiber/v2"
)

func GetUserSelf(c *fiber.Ctx) error {
	var user = model.User{}

	if err := user.GetSelf(c); err != nil {
		return err
	}

	return response.Data(c, &user)
}

func UpdateUserSelf(c *fiber.Ctx) error {
	var user model.User

	if err := user.ScanBody(c); err != nil {
		return err
	}

	if err := user.UpdateSelf(c); err != nil {
		return err
	}

	return response.RecordUpdated(c, user)
}
