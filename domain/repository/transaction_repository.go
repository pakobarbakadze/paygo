package repository

import (
	"paygo/domain/model"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Create(tx *gorm.DB, transaction *model.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *TransactionRepository) CreateLedgerEntry(tx *gorm.DB, entry *model.LedgerEntry) error {
	return tx.Create(entry).Error
}
