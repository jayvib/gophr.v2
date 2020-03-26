package main

import (
  "github.com/gin-gonic/gin"
  "gophr.v2/config"
  "gophr.v2/user/api/v1/http"
  "gophr.v2/user/repository/mysql"
  "gophr.v2/user/repository/mysql/driver"
  "gophr.v2/user/service"
  "log"
)

func main() {
  conf, err := config.New(config.DevelopmentEnv)
  if err != nil {
    log.Fatal(err)
  }
  mysqldb, err := driver.InitializeDriver(conf)
  if err != nil {
    log.Fatal(err)
  }
  defer func() {
    e := mysqldb.Close()
    if e != nil {
      panic(e)
    }
  }()
  repo := mysql.New(mysqldb)
  svc := service.New(repo)
  r := gin.New()
  http.RegisterHandlers(r, svc)
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }

}
