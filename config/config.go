package config

import (
 "github.com/spf13/viper"
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

func New(env Env) (*Config, error) {
	var configPath string
	switch env {
	case DevelopmentEnv:
		configPath = "$HOME/.gophr/testenv/"
	case StageEnv:
		configPath = "$HOME/.gophr/stageenv/"
	case ProdEnv:
		configPath = "$HOME/.gophr/"
	}

	var err error
	once.Do(func(){
		conf, err = loadConfig(
			SetConfigType("yaml"),
			SetConfig("config"),
			AddConfigPath(configPath),
		)
	})
	if err != nil {
		return nil, err
	}
	return conf, nil
}

type Config struct {
	MySQL MySQL `json:"mysql"`
}

type MySQL struct {
	User string
	Password string
	Host    string
	Port string
	Database string
}

func AddConfigPath(path string) func() {
	return func() {
		viper.AddConfigPath(path)
	}
}

func SetConfig(name string) func() {
	return func() {
		viper.SetConfigName(name)
	}
}

func SetConfigType(t string) func() {
	return func(){
		viper.SetConfigType(t)
	}
}

func LoadConfig(opts ...func()) (*Config, error) {
	return loadConfig(opts...)
}

func loadConfig(opts ...func()) (*Config, error) {
	for _, opt := range opts {
		opt()
	}
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	c := new(Config)
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}
	return c, nil
}
