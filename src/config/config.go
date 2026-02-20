package config

type Config struct {
	Version string
}

var AppConfig *Config

func InitializeConfig() {
	AppConfig = &Config{}
}
