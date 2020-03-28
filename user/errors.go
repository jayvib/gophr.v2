package user

import (
  "errors"
  "fmt"
  "strings"
)

var (
	ErrUserNameExists = errors.New("user: can't do operation because user exists")
	ErrEmailExists    = errors.New("user: can't do operation because email exists")

  ErrEmptyUsername      = errors.New("user: cannot process because username is empty")
  ErrEmptyEmail         = errors.New("user: cannot process because email is empty")
  ErrEmptyPassword      = errors.New("user: cannot process because password is empty")
  ErrUserExists         = errors.New("user: cannot create user because it is already exists")
  ErrNotFound           = errors.New("user: item not found")
  ErrInvalidCredentials = errors.New("user: invalid credentials")
)

func NewError(origErr error) *Error {
  e := &Error{
    origErr: origErr,
    context: make(map[interface{}]interface{}),
  }
  e.setMessage()
  return e
}

type Error struct {
  origErr error
  message string
  context map[interface{}]interface{}
}

func (s *Error) Error() string {
  var b strings.Builder
  _, _ = fmt.Fprintf(&b, "%s: %s", s.message, s.origErr)
    for k, v := range s.context {
      _, _ = fmt.Fprintf(&b, " %v: %v", k, v)
    }
  return b.String()
}

func (s *Error) AddContext(k, v interface{}) *Error {
  s.context[k] = v
  return s
}

func (s *Error) Unwrap() error {
  return s.origErr
}

func (s *Error) Message() string {
  return s.message
}

func (s *Error) getMessage() string {
  switch s.origErr {
  case ErrNotFound:
    msg := "Failed getting the user because it didn't exist"

    return msg
  default:
    var b strings.Builder
    _, _ = fmt.Fprint(&b, "Unexpected error: ", s.origErr.Error())
    return b.String()
  }
}

func (s *Error) setMessage() {
  s.message = s.getMessage()
}
