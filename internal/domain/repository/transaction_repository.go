package repository

import (
	"paygo/internal/domain/model"
	"paygo/internal/infra/database"
)

type TransactionRepository struct {
	db database.DB
}

func NewTransactionRepository(db database.DBManager) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) WithTx(tx database.DB) *TransactionRepository {
	return &TransactionRepository{db: tx}
}

func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	return r.db.Create(transaction)
}

func (r *TransactionRepository) CreateLedgerEntry(entry *model.LedgerEntry) error {
	return r.db.Create(entry)
}
