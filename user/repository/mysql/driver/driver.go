package driver

import (
	"database/sql"
	"fmt"
	"gophr.v2/config"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var (
	db *sql.DB
	once sync.Once
)

func InitializeDriver(conf *config.Config) (*sql.DB, error) {
	var err error
	once.Do(func(){
		var e error
		format := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			conf.MySQL.User, conf.MySQL.Password, conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.Database)
		db, e = sql.Open("mysql", format)
		if e != nil {
			err = e
		}
		if e := db.Ping(); e != nil {
			err = e
		}
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
