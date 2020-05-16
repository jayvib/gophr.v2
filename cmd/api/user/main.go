package main

import (
	"github.com/gin-gonic/gin"
	"gophr.v2/config"
	"gophr.v2/user/api/v1/http"
	"gophr.v2/user/repository"
	"gophr.v2/user/service"
)

func main() {
	conf := config.Initialize()

	repo, closer := repository.Get(conf, repository.MySQLRepo)
	defer noOpCloser(closer)

	svc := service.New(repo)

	r := gin.Default()
	http.RegisterHandlers(r, svc)

	if err := r.Run(":4401"); err != nil {
		panic(err)
	}
}

func noOpCloser(c func() error) {
	_ = c()
}