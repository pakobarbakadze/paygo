package database

import (
	"fmt"

	"gorm.io/gorm/clause"
)

func (d *Database) Create(value any) error {
	return d.DB.Create(value).Error
}

func (d *Database) Save(value any) error {
	return d.DB.Save(value).Error
}

func (d *Database) Where(query any, args ...any) DB {
	return &Database{DB: d.DB.Where(query, args...)}
}

func (d *Database) Preload(query string, args ...any) DB {
	return &Database{DB: d.DB.Preload(query, args...)}
}

func (d *Database) First(dest any) error {
	return d.DB.First(dest).Error
}

func (d *Database) Find(dest any) error {
	return d.DB.Find(dest).Error
}

func (d *Database) Clauses(expressions ...clause.Expression) DB {
	return &Database{DB: d.DB.Clauses(expressions...)}
}

func (d *Database) Error() error {
	return d.DB.Error
}

func (d *Database) WithTransaction(fn func(tx DB) error) error {
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
