package handler

import (
	"geophoto/backend/database"
	"geophoto/backend/model"
	"geophoto/backend/utils/response"
	"geophoto/backend/utils/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
)

func getHash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), 14)
	return string(bytes), err
}

func checkHash(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func Login(c *fiber.Ctx) error {
	payload := struct {
		PhoneNumber      string `json:"phone_number" form:"phone_number" validate:"required,number"`
		VerificationCode string `json:"verification_code" form:"verification_code" validate:"required"`
	}{}
	type VerificationCode struct {
		Id         uint   `db:"id"`
		HashedCode string `db:"hashed_code"`
		IsVoided   bool   `db:"is_voided"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return response.InvalidInput(c, "invalid parameter")
	}

	if errors := validate.Struct(payload); errors != nil {
		c.Locals("reason", errors)
		return fiber.NewError(400, "validation failed")
	}

	var vc = VerificationCode{}
	q := "SELECT id, hashed_code, is_voided from verification_code WHERE phone_number=? ORDER BY id DESC"
	if err := database.Cursor.Get(&vc, q, payload.PhoneNumber); err != nil {
		return response.InvalidInput(c, "no available verification code, POST /api/ask/verificatoin to get one")
	}

	if vc.IsVoided || !checkHash(payload.VerificationCode, vc.HashedCode) {
		return response.InvalidInput(c, "invalid verification code")
	}

	// Passed Verification

	// Get User Data
	user := model.User{PhoneNumber: payload.PhoneNumber}
	if err := user.GetByPhoneNumber(); err != nil {
		// Create user if not exist
		if err := user.Create(); err != nil {
			return err
		}
	}

	// Void code
	stmt := "UPDATE verification_code set is_voided=true WHERE id=?"
	if _, err := database.Cursor.Exec(stmt, vc.Id); err != nil {
		return err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.Id
	claims["phone_number"] = user.PhoneNumber
	//claims["exp"] = time.Now().Add(time.Day * 168).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	return response.Data(c, t)
}

func AskVerificationCode(c *fiber.Ctx) error {
	payload := struct {
		PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required,number"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return response.InvalidInput(c, "invalid parameter")
	}

	if errors := validate.Struct(payload); errors != nil {
		c.Locals("reason", errors)
		return fiber.NewError(400, "validation failed")
	}

	length := 6
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(randInt(48, 57))
	}
	verificationCode := string(bytes)
	hashedCode, err := getHash(verificationCode)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO verification_code (phone_number, hashed_code) VALUES (?, ?)"
	if _, err := database.Cursor.Exec(stmt, payload.PhoneNumber, hashedCode); err != nil {
		return err
	}

	// Return the code directly, for development only
	// should send the code via SMS or email in production
	return response.Data(c, fiber.Map{"verification_code": verificationCode})
}
