package database

import (
	"fmt"
	"log"
	"paygo/internal/domain/model"
)

func (d *Database) Migrate() error {
	err := d.DB.AutoMigrate(
		&model.User{},
		&model.Wallet{},
		&model.Transaction{},
		&model.LedgerEntry{},
		&model.Account{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}
