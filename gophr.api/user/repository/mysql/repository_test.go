//+build unit

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jayvib/golog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophr.v2/gophr.api/user"
	"os"
	"testing"
	"time"
)

var debug = flag.Bool("debug", false, "Debug")

var defaultCtx = context.Background()

func TestMain(t *testing.M) {
	flag.Parse()
	if *debug {
		golog.Info("Debug Level")
		golog.SetLevel(golog.DebugLevel)
	}
	os.Exit(t.Run())
}

func setup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *sqlmock.Rows) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	rows := sqlmock.NewRows([]string{
		"id", "userId", "username", "email", "password", "created_at", "updated_at", "deleted_at",
	})
	return db, mock, rows
}

func TestRepository_GetByEmail(t *testing.T) {
	t.Run("Found", func(t *testing.T){
		db, mock, rows := setup(t)
		repo := New(db)
		mockUser := &user.User{
			ID: 1,
			UserID: "testid123",
			Username: "unit.test",
			Email: "unit.test@golang.com",
			Password: "qwerty",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// add the expected output to the rows
		rows.AddRow(
			mockUser.ID,
			mockUser.UserID,
			mockUser.Username,
			mockUser.Email,
			mockUser.Password,
			mockUser.CreatedAt,
			mockUser.UpdatedAt,
			mockUser.DeletedAt,
		)

		query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE email = ?"
		mock.ExpectQuery(query).WillReturnRows(rows)
		u, err := repo.GetByEmail(defaultCtx, "unit.test@golang.com")
		checkErr(t, err)
		assert.Equal(t, mockUser, u)
	})

	t.Run("Not Found", func(t *testing.T){
		db, mock, _ := setup(t)
		repo := New(db)
		query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE email = ?"
		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)
		u, err := repo.GetByEmail(defaultCtx, "unit.test@golang.com")
		assert.Nil(t, u)
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("Unexpected error", func(t *testing.T){
		db, mock, _ := setup(t)
		repo := New(db)
		query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE email = ?"
		mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)
		u, err := repo.GetByEmail(defaultCtx, "unit.test@golang.com")
		assert.Nil(t, u)
		assert.Error(t, err)
		err = errors.Unwrap(err)
		t.Logf("%v\n", err)
	})
}

func TestRepository_GetByID(t *testing.T) {
	db, mock, rows := setup(t)
	repo := New(db)
	mockUser := &user.User{
		ID: 1,
		UserID: "testid123",
		Username: "unit.test",
		Email: "unit.test@golang.com",
		Password: "qwerty",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// add the expected output to the rows
	rows.AddRow(
		mockUser.ID,
		mockUser.UserID,
		mockUser.Username,
		mockUser.Email,
		mockUser.Password,
		mockUser.CreatedAt,
		mockUser.UpdatedAt,
		mockUser.DeletedAt,
	)

	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE id = ?"
	mock.ExpectQuery(query).WillReturnRows(rows)
	u, err := repo.GetByID(defaultCtx, mockUser.ID)
	checkErr(t, err)
	assert.Equal(t, mockUser, u)
}

func TestRepository_GetByUsername(t *testing.T) {
}

func TestRepository_Save(t *testing.T) {
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}
