package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Cursor *sqlx.DB

var TableNames = map[string]string{
	"Photo": "photo",
	"User":  "user",
}
