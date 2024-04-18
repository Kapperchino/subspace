package util

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Debug               bool   `required:"true" default:"false"`
	Port                int    `required:"true" default:"80"`
	JWTSecret           string `required:"true" default:"devSecret" split_words:"true"`
	DatabaseUrl         string `required:"true" split_words:"true"`
	BucketKeySecret     string `required:"true" split_words:"true"`
	BucketKeyId         string `required:"true" split_words:"true"`
	BucketUrl           string `required:"true" split_words:"true"`
	CloudflareAccountId string `required:"true" split_words:"true"`
	BucketName          string `required:"true" split_words:"true"`
	CloudflareApiKey    string `required:"true" split_words:"true"`
	Email               string `required:"true" split_words:"true"`
}

func NewConfig() *Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal().Err(err).Msg("Missing env configs")
	}
	return &c
}
