package handler

import (
	"database/sql"
	"fmt"
	"geophoto/backend/database"
	"geophoto/backend/model"
	"geophoto/backend/utils/response"
	"geophoto/backend/utils/sqlbuilder"
	"github.com/gofiber/fiber/v2"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

var tableNames = database.TableNames

func CreatePhoto(c *fiber.Ctx) error {
	photo := new(model.Photo)

	var file *multipart.FileHeader
	var err error
	if file, err = c.FormFile("file"); err != nil {
		return err
	}
	photo.FileExtension = filepath.Ext(file.Filename)
	if err := photo.ScanBody(c); err != nil {
		return err
	}

	stmt := sqlbuilder.Insert(photo, tableNames["Photo"], "coordinates", "Point(:latitude,:longitude)")

	tx := database.Cursor.MustBegin()
	defer tx.Rollback()
	if rows, err := tx.NamedQuery(stmt, photo); err != nil {
		return err
	} else {
		defer rows.Close()
		rows.Next()
		if err = rows.StructScan(&photo); err != nil {
			log.Println(err)
		}
	}

	if err := c.SaveFile(file, filepath.Join(
		".",
		os.Getenv("IMAGE_PATH"),
		photo.Id+photo.FileExtension,
	)); err != nil {
		return err
	}

	tx.Commit()

	return response.RecordCreated(c, photo)
}

func GetPhoto(c *fiber.Ctx) error {
	var photo model.Photo

	stmt := sqlbuilder.Select(&photo, tableNames["Photo"], "WHERE id=?")
	if err := database.Cursor.Get(&photo, stmt, c.Params("id")); err != nil {
		return err
	}

	return response.Data(c, &photo)
}

func GetAllPhotos(c *fiber.Ctx) error {
	var photos []model.Photo

	stmt := sqlbuilder.Select(&photos, tableNames["Photo"])
	if err := database.Cursor.Select(&photos, stmt); err != nil {
		return err
	}

	return response.Data(c, &photos)
}

func UpdatePhoto(c *fiber.Ctx) error {
	photo := new(model.Photo)

	if err := photo.ScanBody(c); err != nil {
		return err
	}

	var err error
	var file *multipart.FileHeader
	if file, err = c.FormFile("file"); err == nil {
		photo.FileExtension = filepath.Ext(file.Filename)
	}

	stmt := sqlbuilder.Replace(photo, tableNames["Photo"], "coordinates", "Point(:latitude,:longitude)")
	tx := database.Cursor.MustBegin()
	defer tx.Rollback()
	if rows, err := tx.NamedQuery(stmt, photo); err != nil {
		return err
	} else {
		defer rows.Close()
		if !rows.Next() {
			return response.NotFound(c)
		}
		if err = rows.StructScan(&photo); err != nil {
			log.Println(err)
		}
	}

	if file != nil {
		if err := c.SaveFile(file, filepath.Join(
			".",
			os.Getenv("IMAGE_PATH"),
			photo.Id+photo.FileExtension,
		)); err != nil {
			return err
		}
	}

	tx.Commit()

	return response.RecordCreated(c, photo)
}

func DeletePhoto(c *fiber.Ctx) error {
	var err error
	id := c.Params("id")

	var result sql.Result
	if result, err = database.Cursor.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE id=?", tableNames["Photo"]),
		id,
	); err != nil {
		return err
	}

	if rowsDeleted, err := result.RowsAffected(); rowsDeleted == 0 || err != nil {
		return response.NotFound(c)
	}

	return response.RecordDeleted(c)
}
