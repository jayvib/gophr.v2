package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gophr.v2/config"
	"gophr.v2/user/service"
	"gophr.v2/view"
	"log"

	mysqldriver "gophr.v2/driver/mysql"
	imagemysql "gophr.v2/image/repository/mysql"
	imageservice "gophr.v2/image/service"
	sessionfilerepo "gophr.v2/session/repository/file"
	sessionservice "gophr.v2/session/service"
	usermysql "gophr.v2/user/repository/mysql"
)

var (
	conf *config.Config
)

func init() {
	flag.Parse()
	conf = config.Load()
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

	userRepo := usermysql.New(driver)
	userService := service.New(userRepo)

	sessionRepo := sessionfilerepo.New("./sessions.json")
	sessionService := sessionservice.New(sessionRepo)
	r := gin.Default()

	imageRepo := imagemysql.New(driver)
	fs := afero.NewOsFs()
	imageService := imageservice.New(imageRepo, fs, nil)

	view.RegisterRoutes(r, userService, sessionService, imageService,
		"v2/templates/**/*.html",
		"v2/templates/layout.html",
		"v2/assets/",
		"v2/data/images/")

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
