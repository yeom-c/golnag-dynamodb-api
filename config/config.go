package config

import (
	"strconv"

	"github.com/yeom-c/golnag-dynamodb-api/util/env"
)

type Config struct {
	Port        int
	Timeout     int
	Dialect     string
	DatabaseURI string
}

func GetConfig() *Config {
	return &Config{
		Port:        parseEnvToInt("PORT", "8080"),
		Timeout:     parseEnvToInt("TIMEOUT", "5"),
		Dialect:     env.GetEnv("DIALECT", "sqlite3"),
		DatabaseURI: env.GetEnv("DATABASE_URI", ":memory:"),
	}
}

func parseEnvToInt(key, defaultValue string) int {
	value, err := strconv.Atoi(env.GetEnv(key, defaultValue))
	if err != nil {
		return 0
	}

	return value
}
