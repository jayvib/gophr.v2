package gophr

import "errors"

var (
	ErrorNotFound           = errors.New("item not found")
	ErrorSessionExpired     = errors.New("session expired")
	ErrorInvalidCredentials = errors.New("invalid credentials")
)
