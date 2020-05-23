package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	log "github.com/jayvib/golog"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gophr.v2/user"
	"gophr.v2/user/userutil"
	"time"
)

var _ user.Repository = (*Repository)(nil)

func New(conn *sql.DB) *Repository {
	return &Repository{conn: conn}
}

type Repository struct {
	conn *sql.DB
}

func (r *Repository) GetByUserID(ctx context.Context, userID string) (u *user.User, err error) {
	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE userId = ?"
	return r.doQuerySingleReturn(ctx, query, userID)
}

func (r *Repository) GetByID(ctx context.Context, id interface{}) (u *user.User, err error) {
	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE id = ?"
	return r.doQuerySingleReturn(ctx, query, id)
}
func (r *Repository) GetByEmail(ctx context.Context, email string) (u *user.User, err error) {
	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE email = ?"
	return r.doQuerySingleReturn(ctx, query, email)
}
func (r *Repository) GetByUsername(ctx context.Context, uname string) (*user.User, error) {
	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE username = ?"
	return r.doQuerySingleReturn(ctx, query, uname)
}
func (r *Repository) Save(ctx context.Context, usr *user.User) (err error) {
	query := "INSERT INTO user(userId, username, email, password, created_at, updated_at) VALUES(?,?,?,?,?,?)"
	return r.doSave(func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, query,
			usr.UserID,
			usr.Username,
			usr.Email,
			usr.Password,
			usr.CreatedAt,
			usr.UpdatedAt,
		)
		if err != nil {
			return r.checkError(err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			return r.checkError(err)
		}

		usr.ID = uint(id)

		return nil
	})
}
func (r *Repository) Update(ctx context.Context, usr *user.User) error {
	query := "UPDATE user SET userId=?, username=?, email=?, password=?, updated_at=? WHERE id=?"
	return r.doSave(func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, query,
			usr.UserID,
			usr.Username,
			usr.Email,
			usr.Password, // TODO: Remove
			usr.UpdatedAt,
			usr.ID,
		)
		return err
	})
}
func (r *Repository) Delete(ctx context.Context, id interface{}) error {
	query := "DELETE FROM user WHERE id = ?"
	return r.doSave(func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected != 1 {
			return errors.New("repository/mysql: number of rows affected is more than 1")
		}
		return nil
	})
}
func (r *Repository) GetAll(ctx context.Context, cursor string, num int) (users []*user.User, nextCursor string, err error) {
	query := `
		SELECT 
			id, userId, username, email, password, created_at, updated_at, deleted_at 
		FROM 
			user 
		WHERE 
			created_at > ? 
		ORDER BY 
			created_at 
		LIMIT ?`

	var decodedCursor time.Time
	if cursor != "" {
		decodedCursor, err = userutil.DecodeCursor(cursor)
		if err != nil {
			return nil, "", err
		}
	} else {
		decodedCursor = time.Now().AddDate(-100, 0, 0)
	}
	log.Debug("cursor:", decodedCursor)

	res, err := r.doQuery(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	// Generate next pagination cursor
	if len(res) == int(num) {
		nextCursor = userutil.EncodeCursor(*res[len(res)-1].CreatedAt)
	}

	log.Debug("Result:", res)
	log.Debug("Next cursor:", nextCursor)
	return res, nextCursor, nil
}

func (r *Repository) doSave(fn func(tx *sql.Tx) error) (err error) {
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

	err = fn(tx)
	if err != nil {
		return
	}

	return nil

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

	if len(users) == 0 {
		return nil, user.ErrNotFound
	}

	return users[0], nil
}
func (r *Repository) checkError(err error) error {
	var cerr error
	switch err {
	case sql.ErrNoRows:
		cerr = user.ErrNotFound
	case nil:
		cerr = nil
	default:
		cerr = fmt.Errorf("mysql: unexpected error %w", err)
	}
	return cerr
}
