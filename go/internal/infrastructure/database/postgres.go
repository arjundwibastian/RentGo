package database

import (
	"errors"
	"fmt"
	"log"
	"milestone-02/internal/config"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
	
    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("database unreachable: %w", err)
    }

	if err := runMigrations(cfg); err != nil {
        return nil, err
    }
    return db, nil
}

func runMigrations(cfg *config.Config) error {
	sourceURL := os.Getenv("MIGRATIONS_PATH")
	if sourceURL == "" {
		sourceURL = "file://migrations"
	}
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		url.QueryEscape(cfg.Database.Password),
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no new migrations")
			return nil
		}
		return fmt.Errorf("migrate up: %w", err)
	}

	log.Println("database migrations applied")
	return nil
}