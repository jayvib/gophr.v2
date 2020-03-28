package main

import (
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/jayvib/golog"
  "gophr.v2/config"
  "gophr.v2/user/api/v1/http"
  "gophr.v2/user/repository/file"
  "gophr.v2/user/service"
  "log"
)

var (
  envF = flag.String("env", "devel", "Environment. [devel/stage/prod]")
  debug = flag.Bool("debug", false, "Debugging")
)
var (
  conf *config.Config
)

func init() {
  flag.Parse()
  initializeConfig()
  if *debug {
    golog.Info("DEBUGGING MODE")
    golog.SetLevel(golog.DebugLevel)
  }
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
  switch *envF {
  case "devel":
    env = config.DevelopmentEnv
  case "stage":
    env = config.StageEnv
  case "prod":
    env = config.ProdEnv
  }
  conf, err = config.New(env)
  if err != nil {
    log.Fatal(err)
  }
}

