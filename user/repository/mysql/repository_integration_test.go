//+build integration,mysql

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
	"gophr.v2/user/repository/mysql/driver"
	"log"
	"os"
	"testing"
)

var debug = flag.Bool("debug", false, "Debug")

var db *sql.DB
var repo user.Repository

func setup() error {
	conf, err := config.New(config.DevelopmentEnv)
	if err != nil {
		return err
	}

	db, err = driver.InitializeDriver(conf)
	if err  != nil {
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

	t.Run("found", func(t *testing.T){
		email := "luffy.monkey@gmail.com"
		want := &user.User{
			ID: 1,
			UserID: "abc123defe34f334df232dsdfweffewe2fecswf",
			Username: "luffy.monkey",
			Email: "luffy.monkey@gmail.com",
			Password: "secretpass",
		}

		got, err := repo.GetByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("not found", func(t *testing.T){
		_, err := repo.GetByEmail(context.Background(), "not.found@gmail.com")
		assert.Error(t, err)
		assert.Equal(t, mysql.ErrNotFound, err)
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
		ID: input.ID,
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

func assertUpdate(t *testing.T, want *user.User, id interface{}) {
	t.Helper()
	got, err := repo.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func setupUpdate(t *testing.T, input *user.User) (teardown func()){
	t.Helper()
	err := repo.Save(context.Background(), input)
	assert.NoError(t, err)
	return func() {
		err = repo.Delete(context.Background(), input.ID)
		assert.NoError(t, err)
	}
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

	err := repo.Save(context.Background(), want)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), want.ID)
	assert.NoError(t, err)

	// check
	_, err = repo.GetByID(context.Background(), want.ID)
	assert.Error(t, err)
	assert.Equal(t, mysql.ErrNotFound, err)
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
	golog.Debug("ID:", want.ID)
	got, err := repo.GetByID(context.Background(), want.ID)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestRepository_GetAll(t *testing.T) {

}
