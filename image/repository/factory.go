package repository

import (
	"gophr.v2/config"
	mysqldriver "gophr.v2/driver/mysql"
	"gophr.v2/image"
	"gophr.v2/image/repository/mysql"
)

type RepoType int

const (
	MySQLRepo RepoType = iota
)

func Get(conf *config.Config, rt RepoType) (image.Repository, func() error) {
	switch rt {
	case MySQLRepo:
		db, err := mysqldriver.Initialize(conf)
		if err != nil {
			panic(err)
		}
		return mysql.New(db), db.Close
	default:
		panic("unknown repository implementation type")
	}
}
