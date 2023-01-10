package model

import "time"

type Photo struct {
	Id          uint      `db:"id" json:"id"`
	UserId      uint      `db:"user_id" json:"user_id"`
	UUID        string    `db:"uuid" json:"uuid"`
	Description *string   `db:"description" json:"description"`
	Timestamp   time.Time `db:"timestamp" json:"timestamp"`
	Latitude    float64   `db:"X(coordinates)" json:"latitude"`
	Longitude   float64   `db:"Y(coordinates)" json:"longitude"`
}
