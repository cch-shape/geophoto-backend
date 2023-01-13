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
	"github.com/nfnt/resize"
	"github.com/valyala/fasthttp"
	"image"
	"image/jpeg"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type Photo struct {
	//Id          string  `db:"id" db_prop:"auto" json:"id"`
	UUID           string  `db:"uuid" db_prop:"key" json:"uuid"`
	UserId         uint    `db:"user_id" json:"user_id"`
	FileName       string  `db:"filename" db_prop:"auto" json:"filename"`
	PhotoUrl       string  `db_cal:"CONCAT('${IMAGE_PATH}/',uuid,'/',filename)" json:"photo_url"`
	ThumbnailUrl1x string  `db_cal:"CONCAT('${IMAGE_PATH}/',uuid,'/1x/',filename)" json:"thumbnail_url_1x"`
	ThumbnailUrl2x string  `db_cal:"CONCAT('${IMAGE_PATH}/',uuid,'/2x/',filename)" json:"thumbnail_url_2x"`
	Description    *string `db:"description" json:"description"`
	Address        *string `db:"address" json:"address"`
	AddressName    *string `db:"address_name" json:"address_name" form:"address_name"`
	Latitude       float64 `db_cal:"X(coordinates)" json:"latitude" validate:"required,number"`
	Longitude      float64 `db_cal:"Y(coordinates)" json:"longitude" validate:"required,number"`
	Timestamp      string  `db:"timestamp" json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05Z"`
}

type Photos []Photo

func (photo *Photo) saveFile(fh *multipart.FileHeader) error {
	// Create directory
	dir := filepath.Join(".", os.Getenv("IMAGE_PATH"), photo.UUID)
	thumbnailDir := filepath.Join(dir, "1x")
	thumbnailDir2 := filepath.Join(dir, "2x")

	// Save resized image (width=200)
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	thumbnail := resize.Resize(512, 0, img, resize.Lanczos3)
	thumbnail2 := resize.Resize(1024, 0, img, resize.Lanczos3)
	os.MkdirAll(thumbnailDir, os.ModePerm)
	os.MkdirAll(thumbnailDir2, os.ModePerm)
	out, err := os.Create(filepath.Join(thumbnailDir, photo.FileName))
	if err != nil {
		return err
	}
	defer out.Close()
	out2, err := os.Create(filepath.Join(thumbnailDir2, photo.FileName))
	if err != nil {
		return err
	}
	defer out2.Close()
	jpeg.Encode(out, thumbnail, nil)
	jpeg.Encode(out2, thumbnail2, nil)

	// Save original image
	return fasthttp.SaveMultipartFile(fh, filepath.Join(dir, photo.FileName))
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
	photo.UserId = *utils.ExtractJwtUserId(c)

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
var photoCreateStmt = sqlbuilder.Insert(
	Photo{},
	database.TableNames["Photo"],
	"coordinates, filename",
	"Point(:latitude,:longitude), :filename",
)

func (photo *Photo) Create(tx *sqlx.Tx, fh *multipart.FileHeader) error {
	photo.FileName = fh.Filename

	if rows, err := tx.NamedQuery(photoCreateStmt, photo); err != nil {
		return err
	} else {
		defer rows.Close()
		rows.Next()
		if err := rows.StructScan(photo); err != nil {
			return err
		}
	}

	if err := photo.saveFile(fh); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

// Get
var photoGetStmt = sqlbuilder.Select(
	Photo{},
	database.TableNames["Photo"],
	"WHERE uuid=? AND user_id=?",
)

func (photo *Photo) Get() error {
	return database.Cursor.Get(photo, photoGetStmt, photo.UUID, photo.UserId)
}

// Select
var photoSelectStmt = sqlbuilder.Select(
	Photo{},
	database.TableNames["Photo"],
	"WHERE user_id=? ORDER BY id DESC",
)

func (photos *Photos) Select(userId uint) error {
	return database.Cursor.Select(photos, photoSelectStmt, userId)
}

// Update
var photoUpdateStmt = sqlbuilder.Update(
	Photo{},
	database.TableNames["Photo"],
	"coordinates=Point(:latitude,:longitude)",
	" AND user_id=:user_id",
)

var photoUpdateFilenameStmt = fmt.Sprintf(
	"UPDATE `%s` SET filename=:filename WHERE uuid=:uuid AND user_id=:user_id",
	database.TableNames["Photo"],
)

//func (photo *Photo) Update(tx *sqlx.Tx, fh *multipart.FileHeader) error {
//	if _, err := tx.NamedExec(photoUpdateStmt, photo); err != nil {
//		return err
//	}
//	if fh != nil {
//		photo.FileName = fh.Filename
//		if _, err := tx.NamedExec(photoUpdateFilenameStmt, photo); err != nil {
//			return err
//		}
//		photo.deleteFile()
//		if err := photo.saveFile(fh); err != nil {
//			return err
//		}
//	}
//
//	tx.Commit()
//
//	if err := photo.Get(); err != nil {
//		return err
//	}
//
//	return nil
//}

// Delete
var photoDeleteStmt = fmt.Sprintf("DELETE FROM `%s` WHERE uuid=? AND user_id=?", database.TableNames["Photo"])

func (photo *Photo) Delete() (sql.Result, error) {
	if result, err := database.Cursor.Exec(photoDeleteStmt, photo.UUID, photo.UserId); err != nil {
		return nil, err
	} else {
		photo.deleteFile()
		return result, nil
	}
}
