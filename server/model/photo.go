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

// ScanBody /* validate and map request data to model */
func (photo *Photo) ScanBody(c *fiber.Ctx) error {
	var err error
	// Parser body
	if err = c.BodyParser(photo); err != nil {
		return err
	}
	if uuid := c.Params("uuid"); len(uuid) != 0 {
		photo.UUID = uuid
	}
	photo.UserId = 0 /* Read user_id from jwt, to be completed */

	// Validation
	if errors := validate.Struct(photo); errors != nil {
		c.Locals("reason", errors)
		return fiber.NewError(400, "validation failed")
	}

	// Transform
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

func (photo *Photo) Create(tx *sqlx.Tx) (*sqlx.Rows, error) {
	return tx.NamedQuery(createStmt, photo)
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

func (photo *Photo) Update(tx *sqlx.Tx) (sql.Result, error) {
	return tx.NamedExec(updateStmt, photo)
}

var updateFilenameStmt = fmt.Sprintf(
	"UPDATE `%s` SET filename=:filename WHERE uuid=:uuid",
	database.TableNames["Photo"],
)

func (photo *Photo) UpdateFilename(tx *sqlx.Tx) (sql.Result, error) {
	return tx.NamedExec(updateFilenameStmt, photo)
}

// Delete
var deleteStmt = fmt.Sprintf("DELETE FROM `%s` WHERE uuid=?", database.TableNames["Photo"])

func (photo *Photo) Delete() (sql.Result, error) {
	return database.Cursor.Exec(deleteStmt, photo.UUID)
}
