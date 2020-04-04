package session

import "time"

type Session struct {
  ID string `json:"id,omitempty"`
  UserID string `json:"userId,omitempty"`
  Expiry time.Time `json:"expiry,omitempty"`
}
