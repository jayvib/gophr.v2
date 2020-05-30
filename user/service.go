package user

import (
	"context"
)

//go:generate mockery --name=Service

type Service interface {
	GetterByUserID
	GetByID(ctx context.Context, id interface{}) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, uname string) (*User, error)
	Save(ctx context.Context, user *User) error
	GetAll(ctx context.Context, cursor string, num int) (user []*User, nextCursor string, err error)
	Delete(ctx context.Context, id interface{}) error
	Update(ctx context.Context, user *User) error
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, user *User) error
}

type GetterByUserID interface {
	GetByUserID(ctx context.Context, userID string) (*User, error)
}

type Decorator func(svc Service) Service

func ApplyDecorators(svc Service, d ...Decorator) Service {
	copySvc := svc
	for _, deco := range d {
		copySvc = deco(copySvc)
	}
	return copySvc
}
