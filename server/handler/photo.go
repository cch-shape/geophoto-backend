package handler

import (
	"geophoto/backend/database"
	"geophoto/backend/model"
	"geophoto/backend/utils/response"
	"github.com/gofiber/fiber/v2"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

var mPhoto model.Photo

func CreatePhoto(c *fiber.Ctx) error {
	var photo model.Photo

	var file *multipart.FileHeader
	var err error
	if file, err = c.FormFile("file"); err != nil {
		return err
	}
	if err := (&photo).ScanBody(c); err != nil {
		return err
	}
	photo.FileName = file.Filename

	tx := database.Cursor.MustBegin()
	defer tx.Rollback()
	if rows, err := photo.Create(tx); err != nil {
		return err
	} else {
		defer rows.Close()
		rows.Next()
		if err = rows.StructScan(&photo); err != nil {
			log.Println(err)
		}
	}

	dir := filepath.Join(".", os.Getenv("IMAGE_PATH"), photo.UUID)
	os.MkdirAll(dir, os.ModePerm)
	if err := c.SaveFile(file, filepath.Join(dir, file.Filename)); err != nil {
		return err
	}

	tx.Commit()

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
	if _, err := photo.Update(tx); err != nil {
		return err
	}

	if file != nil {
		photo.FileName = file.Filename
		if _, err := photo.UpdateFilename(tx); err != nil {
			return err
		}
		dir := filepath.Join(".", os.Getenv("IMAGE_PATH"), photo.UUID)
		os.RemoveAll(dir)
		os.MkdirAll(dir, os.ModePerm)
		if err := c.SaveFile(file, filepath.Join(dir, file.Filename)); err != nil {
			return err
		}
	}

	tx.Commit()

	if err := photo.Get(); err != nil {
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

	os.RemoveAll(filepath.Join(
		".",
		os.Getenv("IMAGE_PATH"),
		photo.UUID,
	))
	return response.RecordDeleted(c)
}
