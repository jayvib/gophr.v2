package config


type Builder interface {
	SetConfigType() Builder
	SetConfigName() Builder
	AddConfigPath() Builder
	Get() (*Config, error)
}

func Build(builder Builder) (*Config, error) {
	return builder.AddConfigPath().SetConfigName().SetConfigType().Get()
}

