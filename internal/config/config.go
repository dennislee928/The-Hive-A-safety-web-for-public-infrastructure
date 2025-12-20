package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Aggregation AggregationConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Mode         string // debug, release, test
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
	APIKeyPrefix  string
}

// AggregationConfig holds aggregation configuration
type AggregationConfig struct {
	TimeWindows map[string]time.Duration // zone_id -> window duration
	Weights     map[string]map[string]float64 // zone_id -> source_type -> weight
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Mode:         getEnv("GIN_MODE", "debug"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "erh_safety"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
			JWTExpiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
			APIKeyPrefix:  getEnv("API_KEY_PREFIX", "sk_"),
		},
		Aggregation: AggregationConfig{
			TimeWindows: map[string]time.Duration{
				"Z1": 60 * time.Second,
				"Z2": 30 * time.Second,
				"Z3": 90 * time.Second,
				"Z4": 120 * time.Second,
			},
			Weights: map[string]map[string]float64{
				"Z1": {
					"infrastructure": 0.4,
					"staff":          0.4,
					"crowd":          0.2,
					"emergency":      0.5,
				},
				"Z2": {
					"infrastructure": 0.5,
					"staff":          0.4,
					"crowd":          0.1,
					"emergency":      0.5,
				},
				"Z3": {
					"infrastructure": 0.35,
					"staff":          0.35,
					"crowd":          0.3,
					"emergency":      0.5,
				},
				"Z4": {
					"infrastructure": 0.3,
					"staff":          0.4,
					"crowd":          0.3,
					"emergency":      0.5,
				},
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

