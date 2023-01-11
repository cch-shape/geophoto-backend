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

	stmt := sqlbuilder.Insert(
		&photo,
		tableNames["Photo"],
		"coordinates, filename",
		"Point(:latitude,:longitude), :filename",
	)

	tx := database.Cursor.MustBegin()
	defer tx.Rollback()
	if rows, err := tx.NamedQuery(stmt, &photo); err != nil {
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
	var photo model.Photo
	uuid := c.Params("uuid")

	stmt := sqlbuilder.Select(&photo, tableNames["Photo"], "WHERE uuid=?")
	if err := database.Cursor.Get(&photo, stmt, uuid); err != nil {
		return err
	}

	return response.Data(c, &photo)
}

func GetAllPhotos(c *fiber.Ctx) error {
	var photos []model.Photo

	stmt := sqlbuilder.Select(&photos, tableNames["Photo"], "ORDER BY id DESC")
	if err := database.Cursor.Select(&photos, stmt); err != nil {
		return err
	}

	return response.Data(c, &photos)
}

func UpdatePhoto(c *fiber.Ctx) error {
	var photo model.Photo

	if err := (&photo).ScanBody(c); err != nil {
		return err
	}

	file, _ := c.FormFile("file")

	stmt := sqlbuilder.Update(&photo, tableNames["Photo"], "coordinates=Point(:latitude,:longitude)")
	tx := database.Cursor.MustBegin()
	defer tx.Rollback()
	if _, err := tx.NamedExec(stmt, &photo); err != nil {
		return err
	}

	if file != nil {
		photo.FileName = file.Filename
		stmt = fmt.Sprintf("UPDATE `%s` SET filename=:filename WHERE uuid=:uuid", tableNames["Photo"])
		if _, err := tx.NamedExec(stmt, &photo); err != nil {
			return err
		}
		log.Println(stmt)
		dir := filepath.Join(".", os.Getenv("IMAGE_PATH"), photo.UUID)
		os.RemoveAll(dir)
		os.MkdirAll(dir, os.ModePerm)
		if err := c.SaveFile(file, filepath.Join(dir, file.Filename)); err != nil {
			return err
		}
	}

	tx.Commit()

	stmt = sqlbuilder.Select(&photo, tableNames["Photo"], "WHERE uuid=?")
	if err := database.Cursor.Get(&photo, stmt, photo.UUID); err != nil {
		return err
	}

	return response.RecordUpdated(c, photo)
}

func DeletePhoto(c *fiber.Ctx) error {
	var err error
	uuid := c.Params("uuid")

	var result sql.Result
	if result, err = database.Cursor.Exec(
		fmt.Sprintf("DELETE FROM `%s` WHERE uuid=?", tableNames["Photo"]),
		uuid,
	); err != nil {
		return err
	}

	if rowsDeleted, err := result.RowsAffected(); rowsDeleted == 0 || err != nil {
		return response.NotFound(c)
	}
	os.RemoveAll(filepath.Join(
		".",
		os.Getenv("IMAGE_PATH"),
		uuid,
	))
	return response.RecordDeleted(c)
}
