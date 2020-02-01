package user

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, uname string) (*User, error)
	Save(ctx context.Context, user *User) error
}
