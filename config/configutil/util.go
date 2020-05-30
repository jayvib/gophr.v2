package configutil

import (
	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"github.com/spf13/viper"
	"gophr.v2/config"
	builderviper "gophr.v2/config/builder/viper"
	"os"
)

func Initialize() *config.Config {
	return initializeConfig()
}

func LoadDefault(env config.Env) (*config.Config, error) {
	return config.New(builderviper.NewViperBuilder(env))
}

func initializeConfig() *config.Config {
	var err error
	var env config.Env
	switch os.Getenv("GOPHR_ENV") {
	case "DEV":
		env = config.DevelopmentEnv
	case "STAGE":
		env = config.StageEnv
	case "PROD":
		gin.SetMode(gin.ReleaseMode)
		env = config.ProdEnv
	}

	golog.Debug("Environment:", viper.Get("env"))
	conf, err := LoadDefault(env)
	if err != nil {
		panic(err)
	}
	return conf
}
