package repository

import (
	"paygo/internal/domain/model"
	"paygo/internal/infra/database"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	db database.DB
}

func NewAccountRepository(db database.DBManager) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) WithTx(tx database.DB) *AccountRepository {
	return &AccountRepository{db: tx}
}

func (r *AccountRepository) FindByID(id uuid.UUID, forUpdate bool) (*model.Account, error) {
	var account model.Account
	query := r.db.Where("id = ?", id).Preload("LedgerEntries")

	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if err := query.First(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) Update(account *model.Account) (*model.Account, error) {
	if err := r.db.Save(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) FindAll() ([]model.Account, error) {
	var accounts []model.Account
	if err := r.db.Find(&accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}
