package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gophr.v2/gophr.api/user"
	log "github.com/jayvib/golog"
)

var (
	ErrNotFound = errors.New("mysql: Item not found")
)

func New(conn *sql.DB) *Repository {
	return &Repository{conn: conn}
}

type Repository struct {
	conn *sql.DB
}

func (r *Repository) GetByID(ctx context.Context, id string) (*user.User, error) {
	return nil, nil
}
func (r *Repository) GetByEmail(ctx context.Context, email string) (u *user.User, err error) {
	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE email = ?"
	row, err := r.conn.QueryContext(ctx, query, email)
	if err != nil {
		log.Debug(err)
		cerr := r.checkError(err)
		return nil, cerr
	}
	defer func() {
		// TODO: Handle the error
		if e := row.Close(); e != nil {
			err = e
		}
	}()

	u = new(user.User)
	for row.Next() {
		err = row.Scan(&u.ID, &u.UserID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
		if err != nil {
			return nil, r.checkError(err)
		}
	}
	if err = row.Err(); err != nil {
		log.Debug(err)
		return nil, r.checkError(err)
	}
	return u, nil
}

func (r *Repository) checkError(err error) error {
	var cerr error
	switch err {
	case sql.ErrNoRows:
		cerr = ErrNotFound
	default:
		cerr = fmt.Errorf("mysql: unexpected error %w", err)
	}
	return cerr
}
func (r *Repository) GetByUsername(ctx context.Context, uname string) (*user.User, error) {
	return nil, nil
}
func (r *Repository) Save(ctx context.Context, usr *user.User) error {
	return nil
}
