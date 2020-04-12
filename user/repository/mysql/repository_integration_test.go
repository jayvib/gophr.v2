//+build integration

package mysql_test

import (
	"context"
	"database/sql"
	"flag"
	"github.com/jayvib/golog"
	"github.com/stretchr/testify/assert"
	"gophr.v2/config"
	"gophr.v2/user"
	"gophr.v2/user/repository/mysql"
	mysqldriver "gophr.v2/driver/mysql"
	"gophr.v2/user/userutil"
	"gophr.v2/util/valueutil"
	"log"
	"os"
	"testing"
	"time"
)

var debug = flag.Bool("debug", false, "Debug")

var db *sql.DB
var repo user.Repository

func setup() error {
	conf, err := config.New(config.DevelopmentEnv)
	if err != nil {
		return err
	}

	db, err = mysqldriver.InitializeDriver(conf)
	if err != nil {
		return err
	}

	repo = mysql.New(db)
	return nil
}

func TestMain(t *testing.M) {
	flag.Parse()
	if *debug {
		golog.SetLevel(golog.DebugLevel)
	}
	if err := setup(); err != nil {
		log.Fatal(err)
	}
	code := t.Run()
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func TestRepository_GetByEmail(t *testing.T) {

	t.Run("found", func(t *testing.T) {
		email := "luffy.monkey@gmail.com"
		want := &user.User{
			ID:       1,
			UserID:   "abc123defe34f334df232dsdfweffewe2fecswf",
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
			Password: "secretpass",
		}

		got, err := repo.GetByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByEmail(context.Background(), "not.found@gmail.com")
		assert.Error(t, err)
		assert.Equal(t, user.ErrNotFound, err)
	})
}

func TestRepository_GetByID(t *testing.T) {
	want := &user.User{
		ID:       1,
		UserID:   "abc123defe34f334df232dsdfweffewe2fecswf",
		Username: "luffy.monkey",
		Email:    "luffy.monkey@gmail.com",
		Password: "secretpass",
	}

	got, err := repo.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestRepository_GetByUsername(t *testing.T) {
	want := &user.User{
		ID:       1,
		UserID:   "abc123defe34f334df232dsdfweffewe2fecswf",
		Username: "luffy.monkey",
		Email:    "luffy.monkey@gmail.com",
		Password: "secretpass",
	}

	got, err := repo.GetByUsername(context.Background(), "luffy.monkey")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestRepository_Update(t *testing.T) {
	input := &user.User{
		UserID:   "abc123defe34f334df232dsdfweffewe2fecswf",
		Username: "luffy.monkey",
		Email:    "luffy.monkey@gmail.com",
		Password: "secretpass",
	}

	teardown := setupUpdate(t, input)
	defer teardown()

	want := &user.User{
		ID:       input.ID,
		UserID:   "abc123defe34f334df232dsdfweffewe2fecswf",
		Username: "luffy.monkey",
		Email:    "luffy.monkeys@gmail.com",
		Password: "secretpass",
	}

	// For update
	input.Email = "luffy.monkeys@gmail.com"

	err := repo.Update(context.Background(), input)
	assert.NoError(t, err)

	assertUpdate(t, want, input.ID)
}

func TestRepository_Delete(t *testing.T) {
	want := &user.User{
		UserID:   "abc123defe34f334df232dsdfweffewe2fecswf",
		Username: "sanji.vinsmoke",
		Email:    "sanji.vinsmoke@gmail.com",
		Password: "secretpass",

		// NOTE: Remove the time-based field value
		// because it create problem during asserting values
	}

	setupDelete(t, want)

	err := repo.Delete(context.Background(), want.ID)
	assert.NoError(t, err)

	// check
	assertDelete(t, want.ID)
}

func TestRepository_Save(t *testing.T) {
	want := &user.User{
		UserID:   "abc123defe34f334df232dsdfweffewe2fecswf",
		Username: "sanji.vinsmoke",
		Email:    "sanji.vinsmoke@gmail.com",
		Password: "secretpass",

		// NOTE: Remove the time-based field value
		// because it create problem during asserting values
	}

	err := repo.Save(context.Background(), want)
	assert.NoError(t, err)

	assertSave(t, want)
	deleteSaved(t, want.ID)
}

func TestRepository_GetAll(t *testing.T) {

	// Save inputs
	input := []*user.User{
		{
			UserID:    "abc123defe34f334df232dsdfweffewe2fecswf1",
			Username:  "sanji.vinsmoke",
			Email:     "sanji.vinsmoke@gmail.com",
			Password:  "secretpass",
			CreatedAt: valueutil.TimePointer(time.Now().UTC()),
		},
		{
			UserID:    "abc123defe34f334df232dsdfweffewe2fecswf2",
			Username:  "zoro.roronoa",
			Email:     "zoro.roronoa@gmail.com",
			Password:  "secretpass",
			CreatedAt: valueutil.TimePointer(time.Now().UTC()),
		},
		{
			UserID:    "abc123defe34f334df232dsdfweffewe2fecswf3",
			Username:  "nami.navigator",
			Email:     "nami.navigator@gmail.com",
			Password:  "secretpass",
			CreatedAt: valueutil.TimePointer(time.Now().UTC()),
		},
	}

	teardown := setupGetAll(t, input)
	defer teardown()

	cursor := getTimeCursor(*input[0].CreatedAt)
	got, _, err := repo.GetAll(context.Background(), cursor, 3)
	assert.NoError(t, err)

	assert.Len(t, got, 3)

	golog.Debug(got)

	// Compare the result from the input
	assertGetAll(t, input, got)
}

func deleteSaved(t *testing.T, id interface{}) {
	t.Helper()
	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
}

func getTimeCursor(t time.Time) string {
	subTime := t.Add(-time.Second)
	cursor := userutil.EncodeCursor(subTime)
	golog.Debugf("Cursor: %s\n", cursor)

	// Get all
	return cursor
}

func assertGetAll(t *testing.T, want, got []*user.User) {
	t.Helper()

	removeDateValue := func(ins []*user.User) {
		for _, in := range ins {
			in.CreatedAt = nil
			in.UpdatedAt = nil
		}
	}

	removeDateValue(want)
	removeDateValue(got)

	assert.Equal(t, want, got)
}

func setupGetAll(t *testing.T, input []*user.User) (teardown func()) {
	t.Helper()
	for _, in := range input {
		err := repo.Save(context.Background(), in)
		assert.NoError(t, err)
	}
	return func() {
		for _, in := range input {
			err := repo.Delete(context.Background(), in.ID)
			assert.NoError(t, err)
		}
	}
}

func assertUpdate(t *testing.T, want *user.User, id interface{}) {
	t.Helper()
	got, err := repo.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func setupUpdate(t *testing.T, input *user.User) (teardown func()) {
	t.Helper()
	err := repo.Save(context.Background(), input)
	assert.NoError(t, err)
	return func() {
		err = repo.Delete(context.Background(), input.ID)
		assert.NoError(t, err)
	}
}

func assertDelete(t *testing.T, id interface{}) {
	t.Helper()
	_, err := repo.GetByID(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, user.ErrNotFound, err)
}

func setupDelete(t *testing.T, input *user.User) interface{} {
	t.Helper()
	err := repo.Save(context.Background(), input)
	assert.NoError(t, err)
	return input.ID
}

func assertSave(t *testing.T, want *user.User) {
	t.Helper()
	got, err := repo.GetByID(context.Background(), want.ID)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
