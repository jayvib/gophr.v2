package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gophr.v2/config"
	"gophr.v2/user/api/v1/http"
	"gophr.v2/user/repository"
	"gophr.v2/user/service"
)

var port = flag.String("port", "4401", "Port")

func main() {
	flag.Parse()
	conf := config.Initialize()

	repo, closer := repository.Get(conf, repository.MySQLRepo)
	defer noOpCloser(closer)

	svc := service.New(repo)

	r := gin.Default()
	http.RegisterHandlers(r, svc)

	if err := r.Run(fmt.Sprintf(":%s", *port)); err != nil {
		panic(err)
	}
}

func noOpCloser(c func() error) {
	_ = c()
}