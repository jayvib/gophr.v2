package configutil

import (
	"github.com/jayvib/golog"
	"github.com/spf13/viper"
	"gophr.v2/config"
	builderviper "gophr.v2/config/builder/viper"
)

func Initialize() *config.Config {
	return initializeConfig()
}

func LoadDefault() (*config.Config, error) {
	return config.New(builderviper.NewViperBuilder())
}

func initializeConfig() *config.Config {
	golog.Debug("Environment:", viper.Get("env"))
	conf, err := LoadDefault()
	if err != nil {
		panic(err)
	}
	return conf
}
