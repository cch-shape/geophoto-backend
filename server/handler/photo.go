package handler

import (
	"database/sql"
	"geophoto/backend/database"
	"geophoto/backend/model"
	"geophoto/backend/utils/response"
	"geophoto/backend/utils/sqlbuilder"
	"github.com/gofiber/fiber/v2"
)

var tableNames = database.TableNames

func CreatePhoto(c *fiber.Ctx) error {
	photo := model.Photo{}

	if err := c.BodyParser(photo); err != nil {
		return err
	}

	return response.PlaceHolder(c)
}

func GetPhoto(c *fiber.Ctx) error {
	var photo model.Photo
	stmt := sqlbuilder.Select(&photo, tableNames["Photo"], "WHERE uuid=?")

	if err := database.Cursor.Get(&photo, stmt, c.Params("uuid")); err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c)
		}
		return response.InternalServerError(c, &err)
	}

	return response.Data(c, &photo)
}

func GetAllPhotos(c *fiber.Ctx) error {
	var photos []model.Photo
	stmt := sqlbuilder.Select(&photos, tableNames["Photo"])

	if err := database.Cursor.Select(&photos, stmt); err != nil {
		return response.InternalServerError(c, &err)
	}

	return response.Data(c, &photos)
}

func UpdatePhoto(c *fiber.Ctx) error {
	return response.PlaceHolder(c)
}

func DeletePhoto(c *fiber.Ctx) error {
	return response.PlaceHolder(c)
}
