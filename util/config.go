package util

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Debug     bool   `required:"true" default:"false"`
	Port      int    `required:"true" default:"3000"`
	JWTSecret string `required:"true" default:"devSecret"`
}

func NewConfig() *Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &c
}
