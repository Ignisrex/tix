package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string

	RedisHost string
	RedisPort string

	ReservationTTLSeconds int
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
		RedisHost: getEnv("REDIS_HOST", "ticket-lock"),
		RedisPort: getEnv("REDIS_PORT", "6379"),
		ReservationTTLSeconds: getEnvInt("RESERVATION_TTL_SECONDS", 180),
	}
}

func (c Config) DBURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func (c Config) RedisAddr() string {
	return c.RedisHost + ":" + c.RedisPort
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}


