package service

import "errors"

var (
  ErrUsernameEmpty = errors.New("user/service: cannot process because username is empty")
  ErrEmptyEmail = errors.New("user/service: cannot process because email is empty")
  ErrEmptyPassword = errors.New("user/service: cannot process because password is empty")
)
