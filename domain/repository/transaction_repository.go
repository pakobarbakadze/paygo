package repository

import (
	"paygo/domain/model"
	"paygo/infra/database"
)

type TransactionRepositoryImpl struct {
	DB database.DBManager
}

func NewTransactionRepository(db database.DBManager) TransactionRepository {
	return &TransactionRepositoryImpl{DB: db}
}

func (r *TransactionRepositoryImpl) Create(tx database.Transaction, transaction *model.Transaction) error {
	return tx.Create(transaction)
}

func (r *TransactionRepositoryImpl) CreateLedgerEntry(tx database.Transaction, entry *model.LedgerEntry) error {
	return tx.Create(entry)
}
