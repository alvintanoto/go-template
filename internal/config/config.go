package config

import (
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv string `envconfig:"ENV" default:"DEVELOPMENT"`
	Port   string `envconfig:"PORT" default:"8080"`
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("reading viper config: %v", err)
	}

	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		envKey := strings.ToUpper(key)

		if _, exists := os.LookupEnv(envKey); !exists {
			os.Setenv(envKey, value)
		}
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Configuration schema validation failed: %v", err)
	}

	return &cfg
}
