package repository

import (
	"paygo/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) FindByID(tx *gorm.DB, id uuid.UUID, forUpdate bool) (*model.Account, error) {
	db := r.DB
	if tx != nil {
		db = tx
	}

	var account model.Account
	query := db.Where("id = ?", id)

	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if err := query.First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) Update(tx *gorm.DB, account *model.Account) (*model.Account, error) {
	db := r.DB
	if tx != nil {
		db = tx
	}

	if err := db.Save(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}
