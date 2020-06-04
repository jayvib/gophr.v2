package viper

import (
	"github.com/jayvib/golog"
	"github.com/spf13/viper"
	"gophr.v2/config"
	"strconv"
)

const (
	defaultConfigType = "yaml"
	defaultConfigPath = "$HOME"
)

type ConfigBuilderOpt func(b *ConfigBuilder)

func SetConfigName(confName string) ConfigBuilderOpt {
	return func(b *ConfigBuilder) {
		b.configName = confName
	}
}

func SetConfigPath(configPath string) ConfigBuilderOpt {
	return func(b *ConfigBuilder) {
		b.configPath = configPath
	}
}

func SetConfigType(configType string) ConfigBuilderOpt {
	return func(b *ConfigBuilder) {
		b.configType = configType
	}
}

func New(env config.Env, opts ...ConfigBuilderOpt) config.Builder {
	initializeViper()
	b := &ConfigBuilder{
		configName: config.GetConfigName(env),
		configPath: defaultConfigPath,
		configType: defaultConfigType,
	}

	for _, opt := range opts {
		opt(b)
	}
	return b
}

type ConfigBuilder struct {
	configName string
	configPath string
	configType string
	// TODO: Put the viper object here
}

func (d *ConfigBuilder) SetConfigType() config.Builder {
	viper.SetConfigType(d.configType)
	return d
}

func (d *ConfigBuilder) SetConfigName() config.Builder {
	viper.SetConfigName(d.configName)
	return d
}

func (d *ConfigBuilder) AddConfigPath() config.Builder {
	viper.AddConfigPath(d.configPath)
	return d
}

func (d *ConfigBuilder) Get() (*config.Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	c := new(config.Config)
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
}

func initializeViper() {
	viper.AutomaticEnv()
	_ = viper.BindEnv()
	viper.SetEnvPrefix("gophr")
	viper.SetDefault("port", "8080")
}

func initializeDebugging() {
	v := viper.Get("debug")
	isDebug, _ := strconv.ParseBool(v.(string))
	if isDebug {
		golog.SetLevel(golog.DebugLevel)
		golog.Warning("GOPHER IS IN DEBUGGING MODE!")
	}
	viper.GetViper()
}
