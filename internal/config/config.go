package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Servidor
	Port string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// SMTP
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	EmailFrom     string
	EmailFromName string
}

func New() *Config {
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return &Config{
		// Server
		Port: getEnv("PORT", "8080"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		// SMTP
		SMTPHost:      getEnv("SMTP_HOST", ""),
		SMTPPort:      smtpPort,
		SMTPUsername:  getEnv("SMTP_USERNAME", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		EmailFrom:     getEnv("EMAIL_FROM", ""),
		EmailFromName: getEnv("EMAIL_FROM_NAME", "redis-go Service"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
