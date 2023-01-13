package handler

import (
	"geophoto/backend/database"
	"geophoto/backend/model"
	"geophoto/backend/utils/response"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
)

var mPhoto model.Photo

func CreatePhoto(c *fiber.Ctx) error {
	var photo model.Photo

	var fh *multipart.FileHeader
	var err error
	if fh, err = c.FormFile("file"); err != nil {
		return err
	}
	if err := (&photo).ScanBody(c); err != nil {
		return err
	}

	tx := database.Cursor.MustBegin()
	defer tx.Rollback()
	if err := photo.Create(tx, fh); err != nil {
		return err
	}

	return response.RecordCreated(c, photo)
}

func GetPhoto(c *fiber.Ctx) error {
	var photo = model.Photo{UUID: c.Params("uuid")}

	if err := photo.Get(); err != nil {
		return err
	}

	return response.Data(c, &photo)
}

func GetAllPhotos(c *fiber.Ctx) error {
	var photos model.Photos

	if err := photos.Select(); err != nil {
		return err
	}

	return response.Data(c, &photos)
}

func UpdatePhoto(c *fiber.Ctx) error {
	var photo model.Photo

	if err := photo.ScanBody(c); err != nil {
		return err
	}

	file, _ := c.FormFile("file")

	tx := database.Cursor.MustBegin()
	defer tx.Rollback()
	if err := photo.Update(tx, file); err != nil {
		return err
	}

	return response.RecordUpdated(c, photo)
}

func DeletePhoto(c *fiber.Ctx) error {
	var photo = model.Photo{UUID: c.Params("uuid")}

	if result, err := photo.Delete(); err != nil {
		return err
	} else {
		if rowsDeleted, err := result.RowsAffected(); rowsDeleted == 0 || err != nil {
			return response.NotFound(c)
		}
	}

	return response.RecordDeleted(c)
}
