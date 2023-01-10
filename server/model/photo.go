package model

import (
	"geophoto/backend/utils"
	"geophoto/backend/utils/validate"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Photo struct {
	Id          string  `db:"id" pk:"1" json:"id"`
	UserId      uint    `db:"user_id" json:"user_id"`
	FileName    string  `db:"filename" json:"filename"`
	PhotoUrl    string  `db_cal:"CONCAT('${IMAGE_PATH}/',id,'/',filename)" json:"photo_url"`
	Description *string `db:"description" json:"description"`
	Latitude    float64 `db_cal:"X(coordinates)" json:"latitude" validate:"required,number"`
	Longitude   float64 `db_cal:"Y(coordinates)" json:"longitude" validate:"required,number"`
	Timestamp   string  `db:"timestamp" json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05.000Z"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func (photo *Photo) ScanBody(c *fiber.Ctx) error {
	var err error
	// Parser body
	if err = c.BodyParser(photo); err != nil {
		return err
	}
	if id := c.Params("id"); len(id) != 0 {
		photo.Id = id
	}
	photo.UserId = 0 /* Read user_id from jwt, to be completed */

	// Validation
	if errors := validate.Struct(photo); errors != nil {
		c.Locals("reason", errors)
		return fiber.NewError(400, "validation failed")
	}

	// Data transform
	var ts *time.Time
	if ts, err = utils.ISO8601StringToTime(photo.Timestamp); err != nil {
		return err
	}
	photo.Timestamp = utils.TimeToMySQLTimeString(ts)

	return nil
}
