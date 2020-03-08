package user

import (
	"context"
)

//go:generate mockery --name=Repository

type Repository interface {
	GetByID(ctx context.Context, id interface{}) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, uname string) (*User, error)
	Save(ctx context.Context, user *User) error
	GetAll(ctx context.Context, page uint) (*User, error)
	Delete(ctx context.Context, id interface{}) error
}
