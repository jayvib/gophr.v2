package mysql

import (
	"context"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gophr.v2/gophr.api/user"
)

type Repository struct {}

func (r *Repository) GetByID(ctx context.Context, id string) (*user.User, error) {
	return nil, nil
}
func (r *Repository) GetByEmail(ctx context.Context, id string) (*user.User, error) {
	return nil, nil
}
func (r *Repository) GetByUsername(ctx context.Context, uname string) (*user.User, error) {
	return nil, nil
}
func (r *Repository) Save(ctx context.Context, usr *user.User) error {
	return nil
}
