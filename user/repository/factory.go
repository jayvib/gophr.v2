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

func GetRepository(conf *config.Config, rt RepoType) user.Repository {
	switch rt {
	case FileRepo:
		return file.New(file.DefaultFileName)
	case MySQLRepo:
		db, err := mysqldriver.InitializeDriver(conf)
		if err != nil {
			panic(err)
		}
		return mysql.New(db)
	default:
		panic("unknown repository implementation type")
	}
}
