package user

import (
	"context"
)

//go:generate mockery --name=Repository

type Repository interface {
	GetByID(ctx context.Context, id interface{}) (*User, error)
	GetByUserID(ctx context.Context, userId string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, uname string) (*User, error)
	Save(ctx context.Context, user *User) error
	GetAll(ctx context.Context, cursor string, num int) (users []*User, nextCursor string, err error)
	Delete(ctx context.Context, id interface{}) error
	Update(ctx context.Context, user *User) error
}
