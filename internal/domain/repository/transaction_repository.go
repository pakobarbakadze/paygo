package repository

import (
	"paygo/internal/domain/model"
	"paygo/internal/infra/database"
)

type TransactionRepository struct {
	DB database.DBManager
}

func NewTransactionRepository(db database.DBManager) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Create(tx database.Transaction, transaction *model.Transaction) error {
	return tx.Create(transaction)
}

func (r *TransactionRepository) CreateLedgerEntry(tx database.Transaction, entry *model.LedgerEntry) error {
	return tx.Create(entry)
}
