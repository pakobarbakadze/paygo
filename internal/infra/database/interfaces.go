package database

import (
	"gorm.io/gorm/clause"
)

type DB interface {
	Create(value any) error
	Save(value any) error
	Where(query any, args ...any) DB
	Preload(query string, args ...any) DB
	First(dest any) error
	Find(dest any) error
	Clauses(clauses ...clause.Expression) DB
	Error() error
}

type DBManager interface {
	DB
	WithTransaction(fn func(tx DB) error) error
	Close() error
	Migrate() error
}
