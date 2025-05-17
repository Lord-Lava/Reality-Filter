package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/reality-filter/internal/core/ports"
)

// Config implements the ports.ConfigProvider interface
type Config struct {
	MongoDB  mongoDBConfig
	Redis    redisConfig
	Postgres postgresConfig
	Log      logConfig
}

type mongoDBConfig struct {
	URI      string
	Database string
}

type redisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type postgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type logConfig struct {
	Level      string
	Format     string
	OutputPath string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	return &Config{
		MongoDB: mongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://admin:password@localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "reality_filter"),
		},
		Redis: redisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Postgres: postgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "admin"),
			Password: getEnv("POSTGRES_PASSWORD", "password"),
			DBName:   getEnv("POSTGRES_DB", "reality_filter"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Log: logConfig{
			Level:      getEnv("LOG_LEVEL", "debug"),
			Format:     getEnv("LOG_FORMAT", "console"),
			OutputPath: getEnv("LOG_OUTPUT_PATH", "stdout"),
		},
	}, nil
}

// Interface implementation methods

func (c *Config) GetMongoDBConfig() ports.MongoDBConfig {
	return &c.MongoDB
}

func (c *Config) GetRedisConfig() ports.RedisConfig {
	return &c.Redis
}

func (c *Config) GetPostgresConfig() ports.PostgresConfig {
	return &c.Postgres
}

func (c *Config) GetLogConfig() ports.LogConfig {
	return &c.Log
}

// MongoDB implementation
func (c *mongoDBConfig) GetURI() string {
	if c.URI != "" {
		return c.URI
	}
	return "mongodb://localhost:27017"
}

func (c *mongoDBConfig) GetDatabase() string {
	return c.Database
}

// Redis implementation
func (c *redisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *redisConfig) GetPassword() string {
	return c.Password
}

func (c *redisConfig) GetDB() int {
	return c.DB
}

// Postgres implementation
func (c *postgresConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// Log implementation
func (c *logConfig) GetLevel() string {
	return c.Level
}

func (c *logConfig) GetFormat() string {
	return c.Format
}

func (c *logConfig) GetOutputPath() string {
	return c.OutputPath
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
