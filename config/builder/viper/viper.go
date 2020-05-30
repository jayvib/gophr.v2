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

type ViperConfigBuilderOpt func(b *ViperConfigBuilder)

func SetViperConfigName(confName string) ViperConfigBuilderOpt {
  return func(b *ViperConfigBuilder) {
    b.configName = confName
  }
}

func SetViperConfigPath(configPath string) ViperConfigBuilderOpt {
  return func(b *ViperConfigBuilder) {
    b.configPath = configPath
  }
}

func SetViperConfigType(configType string) ViperConfigBuilderOpt {
  return func(b *ViperConfigBuilder) {
    b.configType = configType
  }
}

func NewViperBuilder(env config.Env, opts ...ViperConfigBuilderOpt) config.Builder {
  initializeViper()
  b := &ViperConfigBuilder{
    configName: config.GetConfigName(env),
    configPath: defaultConfigPath,
    configType: defaultConfigType,
  }

  for _, opt := range opts {
    opt(b)
  }
  return b
}

type ViperConfigBuilder struct {
  configName string
  configPath string
  configType string
  // TODO: Put the viper object here
}

func (d *ViperConfigBuilder) SetConfigType() config.Builder {
  viper.SetConfigType(d.configType)
  return d
}

func (d *ViperConfigBuilder) SetConfigName() config.Builder {
  viper.SetConfigName(d.configName)
  return d
}

func (d *ViperConfigBuilder) AddConfigPath() config.Builder {
  viper.AddConfigPath(d.configPath)
  return d
}

func (d *ViperConfigBuilder) Get() (*config.Config, error) {
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
