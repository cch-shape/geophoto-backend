package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func ISO8601StringToTime(timestr string) (*time.Time, error) {
	layout := "2006-01-02T15:04:05Z"
	if t, err := time.Parse(layout, timestr); err != nil {
		return nil, err
	} else {
		return &t, nil
	}
}

func TimeToMySQLTimeString(t *time.Time) string {
	layout := "2006-01-02 15:04:05"
	return t.Format(layout)
}

func ExtractJwtUserId(c *fiber.Ctx) *uint {
	if c.Locals("user") != nil {
		jwtPayload := c.Locals("user").(*jwt.Token)
		claims := jwtPayload.Claims.(jwt.MapClaims)
		userId := uint(claims["user_id"].(float64))
		return &userId
	}
	return nil
}
