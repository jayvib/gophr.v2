package main

import (
  "github.com/gin-gonic/gin"
  "gophr.v2/user/api/v1/http"
  "gophr.v2/user/repository/file"
  "gophr.v2/user/service"
  "log"
)

func main() {
  repo := file.New("./db.json")
  svc := service.New(repo)
  r := gin.New()
  http.RegisterHandlers(r, svc)
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}
