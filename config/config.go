package config

import (
	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"sync"
)

type Env int

const (
	DevelopmentEnv = iota
	StageEnv
	ProdEnv
)

const (
	defaultConfigType = "yaml"
	defaultConfigPath = "$HOME"
)

var (
	conf *Config
	once sync.Once
)

func Initialize() *Config {
	initializeViper()
	initializeConfig()
	return conf
}

func New(env Env) (*Config, error) {
	defBuilder := newViperBuilder(env)
	var err error
	once.Do(func() {
		conf, err = build(defBuilder)
	})
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func getConfigName(env Env) string {
	var configName string
	switch env {
	case DevelopmentEnv:
		configName = "config-dev.yaml"
	case StageEnv:
		configName = "config-stage.yaml"
	case ProdEnv:
		configName = "config.yaml"
	}
	return configName
}

type Config struct {
	rwmu sync.RWMutex
	Gophr Gophr `json:"gophr"`
	MySQL MySQL `json:"mysql"`
}

// Clone creates a new address for existing config.
//
// The cloned config is safe to modify upon cloning.
// This is attempt to implement the Prototype Design Pattern
//
// The aim of the Prototype pattern is to have an object or
// a set of objects that is already created at compilation time,
// but which you can clone as many times as you want at runtime.
func (c *Config) Clone() (*Config, error) {
	c.rwmu.RLock()
	defer c.rwmu.RUnlock()
	clonedConfig := new(Config)
	err := copier.Copy(clonedConfig, c)
	if err != nil {
		return nil, err
	}

	return clonedConfig, nil
}

type Gophr struct {
	Port        string `json:"port"`
	Environment string `json:"env"`
	Debug       bool   `json:"debug"`
}

type MySQL struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func initializeConfig() {
	var err error
	var env Env
	switch viper.Get("env") {
	case "DEV":
		env = DevelopmentEnv
	case "STAGE":
		env = StageEnv
	case "PROD":
		gin.SetMode(gin.ReleaseMode)
		env = ProdEnv
	}

	golog.Debug("Environment:", viper.Get("env"))
	_, err = New(env)
	if err != nil {
		panic(err)
	}
}
