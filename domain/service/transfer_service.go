package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransferService struct {
	DB *gorm.DB
}

func NewTransferService(
	db *gorm.DB,
) *TransferService {
	return &TransferService{
		DB: db,
	}
}

func (s *TransferService) TransferMoney(fromAccountID, toAccountID uuid.UUID, amount float64, description string) {
	println("Transfer money service")
}
