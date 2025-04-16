package database

import (
	"gorm.io/gorm/clause"
)

type Transaction interface {
	Create(value any) error
	Save(value any) error
	Where(query any, args ...any) Transaction
	First(dest any) error
	Clauses(clauses ...clause.Expression) Transaction
	Error() error
}

type DBManager interface {
	WithTransaction(fn func(tx Transaction) error) error
	Close() error
	Migrate() error
}
