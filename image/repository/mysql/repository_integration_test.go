package mysql

import (
	"database/sql"
	"gophr.v2/config"
	"gophr.v2/driver/mysql"
	"os"
	"testing"
)

var db *sql.DB

func setup() {
	conf, err := config.New(config.DevelopmentEnv)
	if err != nil {
		panic(err)
	}
	db, err = mysql.InitializeDriver(conf)
	if err != nil {
		panic(err)
	}
}

//const schema = `
//CREATE DATABASE IF NOT EXISTS gophr_test;
//USE gophr_test;
//
//DROP TABLE IF EXISTS images;
//CREATE TABLE images(
//	id int(36) NOT NULL	AUTO_INCREMENT,
//	updated_at datetime DEFAULT NULL,
//	created_at datetime DEFAULT NULL,
//	deleted_at datetime DEFAULT NULL,
//	userId varchar(45) COLLATE utf8_unicode_ci NOT NULL,
//	name varchar(45) COLLATE utf8_unicode_ci NOT NULL,
//	location varchar(45) COLLATE utf8_unicode_ci NOT NULL,
//	description varchar(100) COLLATE utf8_unicode_ci NOT NULL,
//	size int(36) DEFAULT NULL,
// PRIMARY KEY (id)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
//`

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	err := db.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestRepository_Find(t *testing.T) {
}