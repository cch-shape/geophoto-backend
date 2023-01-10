package utils

import (
	"time"
)

func ISO8601StringToTime(timestr string) (*time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
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
