package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Repository struct {
	db *gorm.DB
}

func(r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	return nil, nil
}
func(r *Repository) GetByEmail(ctx context.Context, id string) (*User, error) {
	return nil, nil
}
func(r *Repository) GetByUsername(ctx context.Context, uname string) (*User, error) {
	return nil, nil
}
func(r *Repository) Save(ctx context.Context, user *User) error {
	return nil
}
