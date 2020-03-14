//+build integration,mysql

package mysql_test

import (
	"database/sql"
	"flag"
	"github.com/jayvib/golog"
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
}

func TestRepository_GetByID(t *testing.T) {
}

func TestRepository_GetByUsername(t *testing.T) {
}

func TestRepository_Delete(t *testing.T) {
}

func TestRepository_GetAll(t *testing.T) {
}

func TestRepository_Update(t *testing.T) {
}

func TestRepository_Save(t *testing.T) {
}
