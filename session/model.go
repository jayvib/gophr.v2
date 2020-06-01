package session

import (
	"time"
)

const (
	Duration   = 24 * 3 * time.Hour
	CookieName = "GophrSession"
)

type Session struct {
	ID     string    `json:"id,omitempty"`
	UserID string    `json:"userId,omitempty"`
	Expiry time.Time `json:"expiry,omitempty"`
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
