package session

import (
  "errors"
  "fmt"
)

var (
  ErrNotFound = errors.New("session: session not exists")
)

func NewError(orig error, msg string) error {
  return &Error{origErr: orig, message: msg}
}

type Error struct {
  origErr error
  message string
}

func (e *Error) Error() string {
  return fmt.Sprintf("%s caused by %s\n", e.message, e.origErr)
}

func (e *Error) Message() string {
  return e.message
}
