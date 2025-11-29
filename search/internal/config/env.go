package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

var Envs Config = initConfig()

func initConfig() Config {
	_ = godotenv.Load()

	return Config{
		Port:      getEnv("PORT", "8081"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPassword:getEnv("DB_PASSWORD", "password"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBName:    getEnv("DB_NAME", "tix_db"),
	}
}

func (c Config) DBURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}


