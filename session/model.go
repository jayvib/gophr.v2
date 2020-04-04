package session

import "time"

type Session struct {
  ID string
  UserID string
  Expiry time.Time
}
