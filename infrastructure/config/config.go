package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contiene la configuración de la aplicación
type Config struct {
	Server struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}
	Redis struct {
		Client   string
		Host     string
		Port     string
		Password string
		DB       int
	}
	RateLimit struct {
		Rate  int
		Burst int
	}
	CORS struct {
		AllowOrigins     []string
		AllowMethods     []string
		AllowHeaders     []string
		AllowCredentials bool
		MaxAge           time.Duration
	}
	Cache struct {
		TTL time.Duration
	}

	Api struct {
		Key string
	}
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	config := &Config{}

	config.Api.Key = getEnv("API_PUBLIC_KEY", "")

	// Server config
	config.Server.Port = getEnv("PORT", "8080")
	config.Server.ReadTimeout = time.Duration(getEnvInt("SERVER_READ_TIMEOUT", 15)) * time.Second
	config.Server.WriteTimeout = time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT", 15)) * time.Second
	config.Server.IdleTimeout = time.Duration(getEnvInt("SERVER_IDLE_TIMEOUT", 60)) * time.Second

	// Redis config
	config.Redis.Client = getEnv("REDIS_CLIENT", "redis")
	config.Redis.Host = getEnv("REDIS_HOST", "127.0.0.1")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.DB = getEnvInt("REDIS_DB", 0)

	// Rate limit config
	config.RateLimit.Rate = getEnvInt("RATE_LIMIT", 1)
	config.RateLimit.Burst = getEnvInt("RATE_LIMIT_BURST", 3)

	// CORS config
	config.CORS.AllowOrigins = getEnvSlice("CORS_ALLOW_ORIGINS", "*")
	config.CORS.AllowMethods = getEnvSlice("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS")
	config.CORS.AllowHeaders = getEnvSlice("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization")
	config.CORS.AllowCredentials = getEnvBool("CORS_ALLOW_CREDENTIALS", false)
	config.CORS.MaxAge = time.Duration(getEnvInt("CORS_MAX_AGE", 12)) * time.Hour

	// Cache config
	config.Cache.TTL = time.Duration(getEnvInt("CACHE_TTL_MINUTES", 5)) * time.Minute

	return config
}

func (c *Config) GetRedisURL() string {
	if c.Redis.Password == "" {
		return fmt.Sprintf("%s://%s:%s/%d", c.Redis.Client, c.Redis.Host, c.Redis.Port, c.Redis.DB)
	}
	return fmt.Sprintf("%s://:%s@%s:%s/%d", c.Redis.Client, c.Redis.Password, c.Redis.Host, c.Redis.Port, c.Redis.DB)
}

// getEnv obtiene una variable de entorno o un valor por defecto
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvSlice(key, defaultValue string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return splitAndTrim(value)
	}
	return splitAndTrim(defaultValue)
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}
