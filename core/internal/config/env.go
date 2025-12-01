package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string

	ESHost string
	ESPort string

	SearchServiceURL string
	BookingServiceURL string
}

var Envs Config = initConfig()

func initConfig() Config {

	godotenv.Load()

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "tix_db"),
		ESHost:     getEnv("ES_HOST", "localhost"),
		ESPort:     getEnv("ES_PORT", "9200"),
		SearchServiceURL: getEnv("SEARCH_SERVICE_URL", "http://search:8082"),
		BookingServiceURL: getEnv("BOOKING_SERVICE_URL", "http://booking:8081"),
	}
}

func (c Config) DBURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func (c Config) ESAddresses() []string {
	return []string{"http://" + c.ESHost + ":" + c.ESPort}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}