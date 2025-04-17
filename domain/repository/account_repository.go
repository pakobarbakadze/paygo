package repository

import (
	"paygo/domain/model"
	"paygo/infra/database"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	DB database.DBManager
}

func NewAccountRepository(db database.DBManager) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) FindByID(tx database.Transaction, id uuid.UUID, forUpdate bool) (*model.Account, error) {
	var account model.Account
	query := tx.Where("id = ?", id)

	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if err := query.First(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) Update(tx database.Transaction, account *model.Account) (*model.Account, error) {
	if err := tx.Save(account); err != nil {
		return nil, err
	}
	return account, nil
}
