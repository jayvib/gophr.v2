package config

type Config struct {
	MySQL *MySQL
}

type MySQL struct {
	Username string
	Password string
	Server   string
}
