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

func (r *Repository) GetByID(ctx context.Context, id interface{}) (u *user.User, err error) {
	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE id = ?"
	return r.doQuerySingleReturn(ctx, query, id)
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (u *user.User, err error) {
	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE email = ?"
	return r.doQuerySingleReturn(ctx, query, email)
}

func (r *Repository) doQuerySingleReturn(ctx context.Context, query string, value interface{}) (u *user.User,err error) {
	row, err := r.conn.QueryContext(ctx, query, value)
	if err != nil {
		log.Debug(err)
		return nil, r.checkError(err)
	}
	defer func() {
		if e := row.Close(); err == nil && e != nil {
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
	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE username = ?"
	return r.doQuerySingleReturn(ctx, query, uname)
}
func (r *Repository) Save(ctx context.Context, usr *user.User) error {
	return nil
}
