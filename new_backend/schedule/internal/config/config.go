package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppName string `envconfig:"APP_NAME" default:"my-app"`
	Port    int    `envconfig:"PORT" default:"8080"`
	Debug   bool   `envconfig:"DEBUG" default:"false"`
	DBUrl   string `envconfig:"DB_URL" required:"true"`
}

func Load() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load env config: %v", err)
	}
	return &cfg
}
