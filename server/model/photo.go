package model

import (
	"database/sql"
	"fmt"
	"geophoto/backend/database"
	"geophoto/backend/utils"
	"geophoto/backend/utils/sqlbuilder"
	"geophoto/backend/utils/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type Photo struct {
	//Id          string  `db:"id" db_prop:"auto" json:"id"`
	UUID        string  `db:"uuid" db_prop:"key" json:"uuid"`
	UserId      uint    `db:"user_id" json:"user_id"`
	FileName    string  `db:"filename" db_prop:"auto" json:"filename"`
	PhotoUrl    string  `db_cal:"CONCAT('${IMAGE_PATH}/',uuid,'/',filename)" json:"photo_url"`
	Description *string `db:"description" json:"description"`
	Latitude    float64 `db_cal:"X(coordinates)" json:"latitude" validate:"required,number"`
	Longitude   float64 `db_cal:"Y(coordinates)" json:"longitude" validate:"required,number"`
	Timestamp   string  `db:"timestamp" json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05.000Z"`
}

type Photos []Photo

func (photo *Photo) saveFile(file *multipart.FileHeader) error {
	dir := filepath.Join(".", os.Getenv("IMAGE_PATH"), photo.UUID)
	os.MkdirAll(dir, os.ModePerm)
	return fasthttp.SaveMultipartFile(file, filepath.Join(dir, file.Filename))
}

func (photo *Photo) deleteFile() error {
	return os.RemoveAll(filepath.Join(".", os.Getenv("IMAGE_PATH"), photo.UUID))
}

func (photo *Photo) ScanBody(c *fiber.Ctx) error {
	var err error
	// Map request data to model
	if err = c.BodyParser(photo); err != nil {
		return err
	}
	if uuid := c.Params("uuid"); len(uuid) != 0 {
		photo.UUID = uuid
	}
	photo.UserId = 0 /* Read user_id from jwt, to be completed */

	// Validate data
	if errors := validate.Struct(photo); errors != nil {
		c.Locals("reason", errors)
		return fiber.NewError(400, "validation failed")
	}

	// Transform data
	var ts *time.Time
	if ts, err = utils.ISO8601StringToTime(photo.Timestamp); err != nil {
		return err
	}
	photo.Timestamp = utils.TimeToMySQLTimeString(ts)

	return nil
}

// Create
var createStmt = sqlbuilder.Insert(
	Photo{},
	database.TableNames["Photo"],
	"coordinates, filename",
	"Point(:latitude,:longitude), :filename",
)

func (photo *Photo) Create(tx *sqlx.Tx, file *multipart.FileHeader) error {
	if rows, err := tx.NamedQuery(createStmt, photo); err != nil {
		return err
	} else {
		defer rows.Close()
		rows.Next()
		if err := rows.StructScan(photo); err != nil {
			return err
		}
	}

	if err := photo.saveFile(file); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

// Get
var getStmt = sqlbuilder.Select(
	Photo{},
	database.TableNames["Photo"],
	"WHERE uuid=?",
)

func (photo *Photo) Get() error {
	return database.Cursor.Get(photo, getStmt, photo.UUID)
}

// Select
var selectStmt = sqlbuilder.Select(
	Photo{},
	database.TableNames["Photo"],
	"ORDER BY id DESC",
)

func (photos *Photos) Select() error {
	return database.Cursor.Select(photos, selectStmt)
}

// Update
var updateStmt = sqlbuilder.Update(
	Photo{},
	database.TableNames["Photo"],
	"coordinates=Point(:latitude,:longitude)",
)

var updateFilenameStmt = fmt.Sprintf(
	"UPDATE `%s` SET filename=:filename WHERE uuid=:uuid",
	database.TableNames["Photo"],
)

func (photo *Photo) Update(tx *sqlx.Tx, file *multipart.FileHeader) error {
	if _, err := tx.NamedExec(updateStmt, photo); err != nil {
		return err
	}
	if file != nil {
		photo.FileName = file.Filename
		if _, err := tx.NamedExec(updateFilenameStmt, photo); err != nil {
			return err
		}
		photo.deleteFile()
		if err := photo.saveFile(file); err != nil {
			return err
		}
	}

	tx.Commit()

	if err := photo.Get(); err != nil {
		return err
	}

	return nil
}

// Delete
var deleteStmt = fmt.Sprintf("DELETE FROM `%s` WHERE uuid=?", database.TableNames["Photo"])

func (photo *Photo) Delete() (sql.Result, error) {
	if result, err := database.Cursor.Exec(deleteStmt, photo.UUID); err != nil {
		return nil, err
	} else {
		photo.deleteFile()
		return result, nil
	}
}
