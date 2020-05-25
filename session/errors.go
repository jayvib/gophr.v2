package session

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("session: session not exists")
	ErrItemExists = errors.New("session: session is already exists")
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

func (e *Error) Unwrap() error {
	return e.origErr
}
