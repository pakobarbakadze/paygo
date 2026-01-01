package repository

import (
	"paygo/internal/domain/model"
	"paygo/internal/infra/database"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	DB database.DBManager
}

func NewAccountRepository(db database.DBManager) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) getDB(tx database.Transaction) database.Transaction {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *AccountRepository) FindByID(tx database.Transaction, id uuid.UUID, forUpdate bool) (*model.Account, error) {
	var account model.Account
	query := r.getDB(tx).Where("id = ?", id).Preload("LedgerEntries")

	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if err := query.First(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) Update(tx database.Transaction, account *model.Account) (*model.Account, error) {
	if err := r.getDB(tx).Save(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) FindAll(tx database.Transaction) ([]model.Account, error) {
	var accounts []model.Account
	if err := r.getDB(tx).Find(&accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}
