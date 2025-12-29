package database

import (
	"context"
	"fmt"
	"log"
	"paygo/internal/config"
	"time"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(cfg *config.Config, opts ...Option) (*Database, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *gorm.DB

	for attempt := 1; attempt <= options.maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), options.timeout)
		defer cancel()

		sqlDB, err := sql.Open("pgx", dsn)
		if err != nil {
			if attempt < options.maxRetries {
				log.Printf("Failed to open database (attempt %d/%d): %v. Retrying in %v...",
					attempt, options.maxRetries, err, options.retryDelay)
				time.Sleep(options.retryDelay)
				continue
			}
			return nil, fmt.Errorf("failed to open database after %d attempts: %w", options.maxRetries, err)
		}

		if err := sqlDB.PingContext(ctx); err != nil {
			sqlDB.Close()
			if attempt < options.maxRetries {
				log.Printf("Failed to ping database (attempt %d/%d): %v. Retrying in %v...",
					attempt, options.maxRetries, err, options.retryDelay)
				time.Sleep(options.retryDelay)
				continue
			}
			return nil, fmt.Errorf("failed to ping database after %d attempts: %w", options.maxRetries, err)
		}

		sqlDB.SetMaxIdleConns(options.maxIdleConns)
		sqlDB.SetMaxOpenConns(options.maxOpenConns)
		sqlDB.SetConnMaxLifetime(options.connMaxLifetime)

		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn:                 sqlDB,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err != nil {
			sqlDB.Close()
			return nil, fmt.Errorf("failed to initialize GORM: %w", err)
		}

		log.Println("Database connection established")
		return &Database{DB: db}, nil
	}

	return nil, fmt.Errorf("unexpected error in database connection")
}

func Setup(cfg *config.Config, opts ...Option) (*Database, error) {
	db, err := NewDatabase(cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Migrate(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	return sqlDB.Close()
}
