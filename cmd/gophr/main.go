package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gophr.v2/config"
	imagerepo "gophr.v2/image/repository"
	sessionrepo "gophr.v2/session/repository"
	userrepo "gophr.v2/user/repository"
	"gophr.v2/user/service"
	"gophr.v2/view"
	"gophr.v2/view/middleware"
	"log"

	imageservice "gophr.v2/image/service"
	sessionservice "gophr.v2/session/service"
)

var (
	conf *config.Config
)

func init() {
	flag.Parse()
	conf = config.Initialize()
	initializeDebugging()
}

func main() {

	userRepo, closer := userrepo.Get(conf, userrepo.MySQLRepo)
	defer noOpClose(closer)
	userService := service.New(userRepo)

	sessionRepo := sessionrepo.Get(sessionrepo.FileRepo)
	sessionService := sessionservice.New(sessionRepo)

	imageRepo, closer := imagerepo.Get(conf, imagerepo.MySQLRepo)
	defer noOpClose(closer)

	fs := afero.NewOsFs()
	imageService := imageservice.New(imageRepo, fs, nil)

	r := gin.Default()
	v1Routers := r.Group("/v1")
	securedRouter := v1Routers.Use(middleware.RequireLogin(sessionService))

	view.RegisterRoutes(r, securedRouter, userService, sessionService, imageService,
		"v2/templates/**/*.html",
		"v2/templates/layout.html",
		"v2/assets/",
		"data/images/")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func initializeDebugging() {
	if viper.GetBool("debug") {
		golog.Info("DEBUGGING MODE")
		golog.SetLevel(golog.DebugLevel)
	}
}

func noOpClose(fn func() error) {
	_ = fn()
}
