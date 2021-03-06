package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jayvib/golog"
	"gophr.v2/config"
	"net/url"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

func Initialize(conf *config.Config) (*sql.DB, error) {
	var err error

	// Using Singleton Pattern
	once.Do(func() {
		var e error
		format := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			conf.MySQL.User, conf.MySQL.Password, conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.Database)
		golog.Info(format)
		val := url.Values{}
		val.Add("parseTime", "1")
		val.Add("loc", "Asia/Manila")
		dsn := fmt.Sprintf("%s?%s", format, val.Encode())
		golog.Debug("DSN:", dsn)
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
