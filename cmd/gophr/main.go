package main

import (
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/jayvib/golog"
  "github.com/spf13/viper"
  "gophr.v2/config"
  "gophr.v2/user/repository/file"
  "gophr.v2/user/service"
  "gophr.v2/view"
  "log"

  sessionfilerepo "gophr.v2/session/repository/file"
  sessionservice "gophr.v2/session/service"
)

var (
  conf *config.Config
)

func init() {
  flag.Parse()
  initializeViper()
  initializeConfig()
  initializeDebugging()
}

func main() {
  userrepo := file.New("./db.json")
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
  if conf.Gophr.Debug {
    golog.Info("DEBUGGING MODE")
    golog.SetLevel(golog.DebugLevel)
  }
}

func initializeViper() {
  viper.AutomaticEnv()
  viper.SetEnvPrefix("gophr")
  viper.SetDefault("port", "8080")
}
