package user

import (
	e "errors"
	"gophr.v2/errors"
)

var (
	ErrNotFound       = errors.ErrorNotFound
	ErrUserNameExists = e.New("user: can't do operation because user exists")
	ErrEmailExists    = e.New("user: can't do operation because email exists")
)
