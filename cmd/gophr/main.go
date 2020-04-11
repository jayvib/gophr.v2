package main

import (
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/jayvib/golog"
  "github.com/spf13/viper"
  "gophr.v2/config"
  "gophr.v2/user/service"
  "gophr.v2/view"
  "log"

  sessionfilerepo "gophr.v2/session/repository/file"
  sessionservice "gophr.v2/session/service"
  usermysql "gophr.v2/user/repository/mysql"
  mysqldriver "gophr.v2/driver/mysql"
)

var (
  conf *config.Config
)

func init() {
  flag.Parse()
  initializeConfig()
  initializeDebugging()
}

func main() {
  driver, err := mysqldriver.InitializeDriver(conf)
  if err != nil {
    log.Fatal(err)
  }
  defer func() {
    err = driver.Close()
    if err != nil {
      golog.Error(err)
    }
  }()

	userrepo := usermysql.New(driver)
  usersvc := service.New(userrepo)

  sessionRepo := sessionfilerepo.New("./sessions.json")
  sessionService := sessionservice.New(sessionRepo)
  r := gin.New()

  view.RegisterRoutes(r, usersvc, sessionService,
    "templates/**/*.html",
    "templates/layout.html")

  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}

func initializeConfig() {
  var err error
  var env config.Env
  switch viper.Get("env") {
  case "DEV":
    env = config.DevelopmentEnv
  case "STAGE":
    env = config.StageEnv
  case "PROD":
    gin.SetMode(gin.ReleaseMode)
    env = config.ProdEnv
  }
  conf, err = config.New(env)
  if err != nil {
    log.Fatal(err)
  }
}

func initializeDebugging() {
  if viper.GetBool("debug") {
    golog.Info("DEBUGGING MODE")
    golog.SetLevel(golog.DebugLevel)
  }
}

