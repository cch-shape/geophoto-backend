package model

import (
	"geophoto/backend/database"
	"geophoto/backend/utils/sqlbuilder"
	"geophoto/backend/utils/validate"
	"github.com/gofiber/fiber/v2"
)

// User struct
type User struct {
	Id           uint   `db:"id" db_prop:"key" json:"id"`
	PhoneNumber  string `db:"phone_number" json:"phone_number" form:"phone_number" validate:"number"`
	Name         string `db:"name" json:"name"`
	ThumbnailURL string `db:"thumbnail_url" json:"thumbnail_url" form:"thumbnail_url"`
	JWT          string
}

func (user *User) ScanBody(c *fiber.Ctx) error {
	var err error
	// Map request data to model
	if err = c.BodyParser(user); err != nil {
		return err
	}

	// Validate data
	if errors := validate.Struct(user); errors != nil {
		c.Locals("reason", errors)
		return fiber.NewError(400, "validation failed")
	}

	return nil
}

// Create
var userCreateStmt = sqlbuilder.Insert(
	User{},
	database.TableNames["User"],
)

func (user *User) Create() error {
	if rows, err := database.Cursor.NamedQuery(userCreateStmt, user); err != nil {
		return err
	} else {
		defer rows.Close()
		rows.Next()
		if err := rows.StructScan(user); err != nil {
			return err
		}
	}

	return nil
}

// Get
var userGetStmt = sqlbuilder.Select(
	User{},
	database.TableNames["User"],
	"WHERE id=?",
)

func (user *User) Get() error {
	return database.Cursor.Get(user, userGetStmt, user.Id)
}

var userGetByPhoneNumberStmt = sqlbuilder.Select(
	User{},
	database.TableNames["User"],
	"WHERE phone_number=?",
)

func (user *User) GetByPhoneNumber() error {
	return database.Cursor.Get(user, userGetByPhoneNumberStmt, user.PhoneNumber)
}

// Update
var userUpdateStmt = sqlbuilder.Update(
	User{},
	database.TableNames["User"],
)

func (user *User) Update() error {
	if _, err := database.Cursor.NamedExec(userUpdateStmt, user); err != nil {
		return err
	}

	if err := user.Get(); err != nil {
		return err
	}

	return nil
}
