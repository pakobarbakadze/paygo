package repository

import (
	"paygo/domain/model"
	"paygo/infra/database"
)

type TransactionRepository struct {
	DB *database.Database
}

func NewTransactionRepository(db *database.Database) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Create(tx *database.Database, transaction *model.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *TransactionRepository) CreateLedgerEntry(tx *database.Database, entry *model.LedgerEntry) error {
	return tx.Create(entry).Error
}
