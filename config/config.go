package config

import (
	"github.com/jayvib/golog"
	"github.com/jinzhu/copier"
	"sync"
)

type Env int

const (
	DevelopmentEnv = iota
	StageEnv
	ProdEnv
)

var (
	conf *Config
	once sync.Once
)

func New(builder Builder) (*Config, error) {
	var err error
	once.Do(func() {
		conf, err = Build(builder)
		conf.init()
	})
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func GetConfigName(env Env) string {
	return getConfigName(env)
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
	rwmu  sync.RWMutex
	Gophr Gophr `json:"gophr"`
	MySQL MySQL `json:"mysql"`
	Debug bool  `json:"debug"`
}

func (c *Config) init() {
	if c.Debug {
		golog.Warning("Gopher is in debug mode!")
		golog.SetLevel(golog.DebugLevel)
	}
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

