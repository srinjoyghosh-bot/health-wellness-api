// internal/config/config.go
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config structure to hold the application configuration
type Config struct {
	Server   ServerConfig
	Database DBConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port int
	Mode string // e.g., development, production
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration string
}

// Load reads the configuration from a file or environment variables
func Load() (*Config, error) {
	var config Config

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(".")        // Current directory
	viper.AddConfigPath("./config") // Config folder

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values (optional)
	setDefaults()

	// Read the configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found. Using environment variables and defaults.")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}

// Set default values
func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "development")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "mysecretpassword")
	viper.SetDefault("database.dbname", "health_db")
	viper.SetDefault("database.sslmode", "disable")

	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiration", "30d")
}

// GetDSN generates Data Source Name for database connection
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode)
}

// ParseExpiration ParseJWTExpiration converts JWT expiration to time.Duration
func (c *JWTConfig) ParseExpiration() (time.Duration, error) {
	duration, err := time.ParseDuration(strings.TrimRight(c.Expiration, "d") + "d")
	if err != nil {
		return 0, fmt.Errorf("invalid JWT expiration format: %w", err)
	}
	return duration, nil
}
