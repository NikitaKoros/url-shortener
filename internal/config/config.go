package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	BaseURL string
	Host    string
	Port    int
	DB      DBConfig
}

type DBConfig struct {
	DSN      string
	Database string
}

var (
	cfg  Config
	once sync.Once
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Error converting %s to int: %v", key, err)
		return defaultValue
	}
	return value
}

func LoadConfig() Config {
	once.Do(func() {
		cfg = Config{
			BaseURL: getEnv("BASE_URL", "http://localhost:5432"),
			Host:    getEnv("HOST", "0.0.0.0"),
			Port:    getEnvAsInt("PORT", 8080),
			DB: DBConfig{
				DSN: getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=123 port=5432"),

				Database: getEnv("DATABASE_NAME", "123"),
			},
		}
		fmt.Println(cfg.BaseURL, cfg.Host, cfg.Port, cfg.DB.DSN, cfg.DB.Database)
	})

	return cfg
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
