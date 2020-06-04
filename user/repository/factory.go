package repository

import (
	"gophr.v2/config"
	mysqldriver "gophr.v2/driver/mysql"
	"gophr.v2/user"
	"gophr.v2/user/repository/file"
	"gophr.v2/user/repository/mysql"
)

type RepoType int

const (
	FileRepo RepoType = iota
	MySQLRepo
)

func Get(conf *config.Config, rt RepoType) (user.Repository, func() error) {
	switch rt {
	case FileRepo:
		return file.New(file.DefaultFileName), noOpClose
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

func noOpClose() error {
	return nil
}
