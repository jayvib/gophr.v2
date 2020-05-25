package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"gophr.v2/config"
	"gophr.v2/image/api/v1"
	imagerepo "gophr.v2/image/repository"
	imageservice "gophr.v2/image/service"
	userrepo "gophr.v2/user/repository"
	userservice "gophr.v2/user/service"
)

var port = flag.String("port", "4402", "Port")

func main() {
	flag.Parse()
	conf := config.Initialize()

	userRepo, closer := userrepo.Get(conf, userrepo.MySQLRepo)
	defer closer()

	userService := userservice.New(userRepo)

	imageRepo, closer := imagerepo.Get(conf, imagerepo.MySQLRepo)
	defer closer()

	fs := afero.NewOsFs()
	imageService := imageservice.New(imageRepo, fs, nil)

	e := gin.Default()
	v1.RegisterRoutes(e, imageService, userService)

	if err := e.Run(fmt.Sprintf(":%v", *port)); err != nil {
		panic(err)
	}
}

func noOpCloser(c func() error) {
	_ = c()
}
