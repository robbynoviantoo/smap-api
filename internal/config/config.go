package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	AppPort      string
	AppEnv       string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	JWTExpireHours int
}

var App Config

func Load() {
	App = Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "smap"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me"),
		JWTExpireHours: 24,
	}

	log.Printf("[Config] Loaded: env=%s, port=%s, db=%s@%s:%s/%s",
		App.AppEnv, App.AppPort,
		App.DBUser, App.DBHost, App.DBPort, App.DBName,
	)
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
