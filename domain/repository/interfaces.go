package repository

import (
	"paygo/domain/model"
	"paygo/infra/database"

	"github.com/google/uuid"
)

type AccountRepository interface {
	FindByID(tx database.Transaction, id uuid.UUID, forUpdate bool) (*model.Account, error)
	Update(tx database.Transaction, account *model.Account) (*model.Account, error)
}

type TransactionRepository interface {
	Create(tx database.Transaction, transaction *model.Transaction) error
	CreateLedgerEntry(tx database.Transaction, entry *model.LedgerEntry) error
}
