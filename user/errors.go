package user

import (
  "errors"
)

var (
	ErrUserNameExists = errors.New("user: can't do operation because user exists")
	ErrEmailExists    = errors.New("user: can't do operation because email exists")

  ErrEmptyUsername      = errors.New("user: cannot process because username is empty")
  ErrEmptyEmail         = errors.New("user: cannot process because email is empty")
  ErrEmptyPassword      = errors.New("user: cannot process because password is empty")
  ErrUserExists         = errors.New("user: cannot create user because it is already exists")
  ErrNotFound           = errors.New("user: user not found")
  ErrInvalidCredentials = errors.New("user: invalid credentials")
)
