package driver

import (
	"database/sql"
	"fmt"
	"gophr.v2/config"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
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

		val := url.Values{}
		val.Add("parseTime", "1")
		val.Add("loc", "Asia/Manila")
		dsn := fmt.Sprintf("%s?%s", format, val.Encode())
		db, e = sql.Open("mysql", dsn)
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
