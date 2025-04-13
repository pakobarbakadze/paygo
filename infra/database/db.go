package database

import (
	"context"
	"fmt"
	"log"
	"paygo/config"
	"paygo/domain/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Add connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Do a ping to verify the connection is valid
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")

	return &Database{DB: db}, nil
}

func Setup(cfg *config.Config) (*Database, error) {
	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Migrate(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func (d *Database) WithTransaction(fn func(tx *Database) error) error {
	gormTx := d.DB.Begin()
	if gormTx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", gormTx.Error)
	}

	tx := &Database{DB: gormTx}

	defer func() {
		if r := recover(); r != nil {
			gormTx.Rollback()
			fmt.Println("Transaction rolled back due to panic:", r)
		}
	}()

	if err := fn(tx); err != nil {
		gormTx.Rollback()
		fmt.Println("Transaction rolled back due to error:", err)
		return err
	}

	return gormTx.Commit().Error
}

func (d *Database) Migrate() error {
	err := d.DB.AutoMigrate(&model.User{}, &model.Wallet{}, &model.Transaction{}, &model.LedgerEntry{}, &model.Account{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("Database migration completed")
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	return sqlDB.Close()
}
