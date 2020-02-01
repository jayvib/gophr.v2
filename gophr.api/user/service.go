package user

import (
	"context"
)

//go:generate mockery --name=Service

type Service interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, uname string) (*User, error)
	Save(ctx context.Context, user *User) error
	GetAndComparePassword(ctx context.Context, username, password string) (*User, error)
}
