package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Port string
	AdminKey string
	APIKey string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration 
}


func Load() (*Config, error) {
	cfg := &Config {
		App: AppConfig{
			Port: os.Getenv("PORT"),
			AdminKey: os.Getenv("ADMIN_KEY"),
			APIKey: os.Getenv("API_KEY"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Expiry: getEnvDuration("JWT_EXPIRY_HOURS", 24) * time.Hour,
		},
	}

	 if err := cfg.validate(); err != nil {
        return nil, err
    }
    return cfg, nil

}

func getEnvDuration (key string, fallbackHours int) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return time.Duration(fallbackHours)
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return time.Duration(fallbackHours)
	}
	return time.Duration(n)
}

func (cfg *Config) validate() error {
    required := map[string]string{
        "DB_HOST":     cfg.Database.Host,
        "DB_USER":     cfg.Database.User,
        "DB_PASSWORD": cfg.Database.Password,
        "DB_NAME":     cfg.Database.Name,
        "DB_PORT":     cfg.Database.Port,
        "JWT_SECRET":  cfg.JWT.Secret,
		"ADMIN_KEY": cfg.App.AdminKey,
		"API_KEY": cfg.App.APIKey,
    }
    for key, value := range required {
        if value == "" {
            return fmt.Errorf("missing required environment variable: %s", key)
        }
    }
    return nil
}