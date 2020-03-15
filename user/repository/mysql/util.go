package mysql

import (
	"encoding/base64"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.999Z07:00"

func EncodeCursor(t time.Time) string {
	return encodeCursor(t)
}

func DecodeCursor(encodedTime string) (time.Time, error) {
	return decodeCursor(encodedTime)
}

// encodeCursor encode time t into a base 64 string.
func encodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)
	return base64.StdEncoding.EncodeToString([]byte(timeString))
}

func decodeCursor(encodedTime string) (time.Time, error) {
	var emptyTime time.Time
	timeString, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return emptyTime, err
	}
	t, err := time.Parse(timeFormat, string(timeString))
	if err != nil {
		return emptyTime, err
	}
	return t, nil
}
