package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppName  string `envconfig:"APP_NAME" default:"my-app"`
	Port     int    `envconfig:"PORT" default:"8080"`
	Debug    bool   `envconfig:"DEBUG" default:"false"`
	DBUrl    string `envconfig:"DB_URL" required:"true"`
	Schedule struct {
		Host string `envconfig:"SCHEDULE_HOST" default:"schedule"`
		Port int    `envconfig:"SCHEDULE_PORT"`
	}
	User struct {
		Host string `envconfig:"USER_HOST" default:"user"`
		Port int    `envconfig:"USER_PORT"`
	}
	RabbitMq struct {
		Uri string `envconfig:"RABBITMQ_URL" default:"amqp://guest:guest@rabbitmq:5672/"`
	}
}

func Load() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load env config: %v", err)
	}
	return &cfg
}
