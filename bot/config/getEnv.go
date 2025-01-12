package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	HTTPPort string

	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	LOG_PATH string
}

func Load() Config {
	config := Config{}

	config.HTTPPort = cast.ToString(coalesce("HTTP_PORT", ":8080"))

	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	config.DB_PORT = cast.ToInt(coalesce("DB_PORT", 5432))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "azizbek"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "123"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "bot"))

	config.LOG_PATH = cast.ToString(coalesce("LOG_PATH", "logs/info.log"))

	return config
}

func coalesce(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
