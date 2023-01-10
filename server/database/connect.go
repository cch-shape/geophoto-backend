package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

func Connect() {
	var err error
	Cursor, err = sqlx.Connect("mysql",
		fmt.Sprintf(
			"%s:%s@(%s:%s)/%s?parseTime=true",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Database")
}
