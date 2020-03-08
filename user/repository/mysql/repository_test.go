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
	"gophr.v2/user"
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
	t.Run("Found", func(t *testing.T) {
		db, mock, rows := setup(t)
		repo := New(db)
		mockUser := &user.User{
			ID:        1,
			UserID:    "testid123",
			Username:  "unit.test",
			Email:     "unit.test@golang.com",
			Password:  "qwerty",
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

	t.Run("Not Found", func(t *testing.T) {
		db, mock, _ := setup(t)
		repo := New(db)
		query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE email = ?"
		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)
		u, err := repo.GetByEmail(defaultCtx, "unit.test@golang.com")
		assert.Nil(t, u)
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("Unexpected error", func(t *testing.T) {
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
		ID:        1,
		UserID:    "testid123",
		Username:  "unit.test",
		Email:     "unit.test@golang.com",
		Password:  "qwerty",
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
	db, mock, rows := setup(t)
	repo := New(db)
	mockUser := &user.User{
		ID:        1,
		UserID:    "testid123",
		Username:  "unit.test",
		Email:     "unit.test@golang.com",
		Password:  "qwerty",
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

	query := "SELECT id,userId,username,email,password,created_at,updated_at,deleted_at FROM user WHERE username = ?"
	mock.ExpectQuery(query).WillReturnRows(rows)
	u, err := repo.GetByUsername(defaultCtx, mockUser.Username)
	checkErr(t, err)
	assert.Equal(t, mockUser, u)
}

func TestRepository_Save(t *testing.T) {
	mockUser := &user.User{
		UserID:    "testid123",
		Username:  "unit.test",
		Email:     "unit.test@golang.com",
		Password:  "qwerty",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	want := &(*mockUser)
	db, mock, _ := setup(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO user").WithArgs(
		mockUser.UserID,
		mockUser.Username,
		mockUser.Email,
		mockUser.Password,
		mockUser.CreatedAt,
		mockUser.UpdatedAt,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	repo := New(db)
	err := repo.Save(context.Background(), mockUser)
	require.NoError(t, err)
	assert.Equal(t, want, mockUser)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update(t *testing.T) {
	mockUser := &user.User{
		UserID:    "testid123",
		Username:  "unit.test",
		Email:     "unit.test@golang.com",
		Password:  "qwerty",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db, mock, _ := setup(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE user").WithArgs(
		mockUser.UserID,
		mockUser.Username,
		mockUser.Email,
		mockUser.Password,
		mockUser.UpdatedAt,
		mockUser.ID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	repo := New(db)
	err := repo.Update(context.Background(), mockUser)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete(t *testing.T) {
	mockUser := &user.User{
		UserID:    "testid123",
		Username:  "unit.test",
		Email:     "unit.test@golang.com",
		Password:  "qwerty",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db, mock, _ := setup(t)
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM user").WithArgs(
		mockUser.ID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	repo := New(db)
	err := repo.Delete(context.Background(), mockUser.ID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAll(t *testing.T) {
	// Add the mock users to the rows
	mockUsers := []*user.User{
		{
			ID: 1,
			UserID:    "testid123",
			Username:  "unit.test",
			Email:     "unit.test@golang.com",
			Password:  "qwerty",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID: 2,
			UserID:    "testid124",
			Username:  "unit.test01",
			Email:     "unit.test01@golang.com",
			Password:  "qwerty",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID: 3,
			UserID:    "testid125",
			Username:  "unit.test02",
			Email:     "unit.test02@golang.com",
			Password:  "qwerty",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	t.Run("All Results", func(t *testing.T){
		db, mock, rows := setup(t)

		// Add the mock users to the row
		for _, u := range mockUsers {
			rows.AddRow(u.ID, u.UserID, u.Username, u.Email, u.Password, u.CreatedAt, u.UpdatedAt, u.DeletedAt)
		}
		// Need to escape the "?" character as per this issue:
		// https://github.com/DATA-DOG/go-sqlmock/issues/70
		query := "SELECT id, userId, username, email, password, created_at, updated_at, deleted_at FROM user WHERE created_at > \\? ORDER BY created_at LIMIT \\?"
		mock.ExpectQuery(query).WillReturnRows(rows)

		repo := New(db)
		cursor := encodeCursor(mockUsers[0].CreatedAt)
		list, nextCursor, err := repo.GetAll(context.Background(), cursor, 3)
		_ = nextCursor
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Len(t, list, 3)
	})
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}
