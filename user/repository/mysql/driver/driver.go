package driver

import (
	"database/sql"
	"fmt"
	"gophr.v2/config"
)

func InitializeDriver(conf *config.Config) (*sql.DB, error) {
	format := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		conf.MySQL.User, conf.MySQL.Password, conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.Name)
	db, err := sql.Open("mysql", format)
	if err != nil {
		return nil, err
	}
	return db, nil
}
