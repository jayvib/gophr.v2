package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	log "github.com/jayvib/golog"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gophr.v2/gophr.api/user"
)

var (
	ErrNotFound = errors.New("repository/mysql: Item not found")
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

func (r *Repository) doQuery(ctx context.Context, query string, args ...interface{}) (users []*user.User, err error) {
	row, err := r.conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Debug(err)
		return nil, r.checkError(err)
	}
	defer func() {
		if e := row.Close(); err == nil && e != nil {
			err = e
		}
	}()

	users = make([]*user.User, 0)
	for row.Next() {
		var u user.User
		err = row.Scan(&u.ID, &u.UserID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
		if err != nil {
			return nil, r.checkError(err)
		}
		users = append(users, &u)
	}
	if err = row.Err(); err != nil {
		log.Debug(err)
		return nil, r.checkError(err)
	}
	return users, nil
}

func (r *Repository) doQuerySingleReturn(ctx context.Context, query string, value interface{}) (u *user.User, err error) {
	users, err := r.doQuery(ctx, query, value)
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (r *Repository) checkError(err error) error {
	var cerr error
	switch err {
	case sql.ErrNoRows:
		cerr = ErrNotFound
	case nil:
		cerr = nil
	default:
		cerr = fmt.Errorf("mysql: unexpected error %w", err)
	}
	return cerr
}
func (r *Repository) GetByUsername(ctx context.Context, uname string) (*user.User, error) {
	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE username = ?"
	return r.doQuerySingleReturn(ctx, query, uname)
}
func (r *Repository) Save(ctx context.Context, usr *user.User) (err error) {
	// When modifying a data, transaction is a good idea
	tx, err := r.conn.Begin()
	if err != nil {
		return r.checkError(err)
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			if e := tx.Rollback(); e != nil {
				// TODO: I think its better to cascade the error
				// For now just log the error.
				log.Error("error while rolling back data in sql:", e)
			}
		}
	}()

	query := "INSERT INTO user(userId, username, email, password, created_at, updated_at) VALUES(?,?,?,?,?,?)"
	_, err = tx.ExecContext(ctx, query,
		usr.UserID,
		usr.Username,
		usr.Email,
		usr.Password,
		usr.CreatedAt,
		usr.UpdatedAt,
	)
	return r.checkError(err)
}

func (r *Repository) Update(ctx context.Context, usr *user.User) error {
	return nil
}

func (r *Repository) Delete(ctx context.Context, id interface{}) error {
	return nil
}

func (r *Repository) GetAll(ctx context.Context, page uint) ([]*user.User, error) {
	return nil, nil
}
