package config

import (
	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	"github.com/spf13/viper"
	"log"
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

func Load() *Config {
	initializeViper()
	initializeConfig()
	return conf
}

func New(env Env) (*Config, error) {

	defBuilder := newDefaultBuilder(env)

	var err error
	once.Do(func() {
		conf, err = build(defBuilder)
		initializeViper()
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

func initializeViper() {
	viper.AutomaticEnv()
	_ = viper.BindEnv()
	viper.SetEnvPrefix("gophr")
	viper.SetDefault("port", "8080")
}

type Config struct {
	Gophr Gophr `json:"gophr"`
	MySQL MySQL `json:"mysql"`
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
	return func() {
		viper.SetConfigType(t)
	}
}

type Builder interface {
	SetConfigType() Builder
	SetConfigName() Builder
	AddConfigPath() Builder
	Get() (*Config, error)
}

func build(builder Builder) (*Config, error) {
	return builder.AddConfigPath().SetConfigName().SetConfigType().Get()
}

func newDefaultBuilder(env Env) Builder {
	return &defaultBuilder{
		configName: getConfigName(env),
		configPath: defaultConfigPath,
		configType: defaultConfigType,
	}
}

type defaultBuilder struct {
	configName string
	configPath string
	configType string
}

func (d *defaultBuilder) SetConfigType() Builder {
	SetConfigType(d.configType)
	return d
}

func (d *defaultBuilder) SetConfigName() Builder {
	SetConfig(d.configName)
	return d
}

func (d *defaultBuilder) AddConfigPath() Builder {
	AddConfigPath(d.configPath)
	return d
}

func (d *defaultBuilder) Get() (*Config, error) {

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	c := new(Config)
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
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
	golog.Info(env)
	_, err = New(env)
	if err != nil {
		log.Fatal(err)
	}
}

