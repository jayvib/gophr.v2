package main

import (
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/jayvib/golog"
  "github.com/spf13/viper"
  "gophr.v2/config"
  "gophr.v2/user/api/v1/http"
  "gophr.v2/user/repository/file"
  "gophr.v2/user/service"
  "log"
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
  repo := file.New("./db.json")
  svc := service.New(repo)
  r := gin.New()
  http.RegisterHandlers(r, svc)
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
}
