package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	RedisURL string `envconfig:"REDIS_URL" default:"localhost:6379"`
}

func NewConfig() *Config {
	var conf Config

	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
