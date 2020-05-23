package config

import (
	"github.com/jayvib/golog"
	"github.com/spf13/viper"
	"strconv"
)

type Builder interface {
	SetConfigType() Builder
	SetConfigName() Builder
	AddConfigPath() Builder
	Get() (*Config, error)
}

func build(builder Builder) (*Config, error) {
	return builder.AddConfigPath().SetConfigName().SetConfigType().Get()
}

func newViperBuilder(env Env) Builder {
	initializeViper()
	return &viperConfigBuilder{
		configName: getConfigName(env),
		configPath: defaultConfigPath,
		configType: defaultConfigType,
	}
}

type viperConfigBuilder struct {
	configName string
	configPath string
	configType string
	// TODO: Put the viper object here
}

func (d *viperConfigBuilder) SetConfigType() Builder {
	viper.SetConfigType(d.configType)
	return d
}

func (d *viperConfigBuilder) SetConfigName() Builder {
	viper.SetConfigName(d.configName)
	return d
}

func (d *viperConfigBuilder) AddConfigPath() Builder {
	viper.AddConfigPath(d.configPath)
	return d
}

func (d *viperConfigBuilder) Get() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	c := new(Config)
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
