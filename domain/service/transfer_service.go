package service

import (
	"errors"
	"fmt"
	"paygo/domain/model"
	"paygo/domain/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransferService struct {
	DB              *gorm.DB
	AccountRepo     *repository.AccountRepository
	TransactionRepo *repository.TransactionRepository
}

func NewTransferService(
	db *gorm.DB,
	accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepository,
) *TransferService {
	return &TransferService{
		DB:              db,
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
	}
}

func (s *TransferService) TransferMoney(fromAccountID, toAccountID uuid.UUID, amount float64, description string) (*model.Transaction, *model.Account, *model.Account, error) {
	var fromAccount, toAccount *model.Account
	var err error

	tx := s.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if fromAccount, err = s.AccountRepo.FindByID(tx, fromAccountID, true); err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	if toAccount, err = s.AccountRepo.FindByID(tx, toAccountID, true); err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	if fromAccount.Status != "active" {
		tx.Rollback()
		return nil, nil, nil, errors.New("source account is not active")
	}

	if toAccount.Status != "active" {
		tx.Rollback()
		return nil, nil, nil, errors.New("destination account is not active")
	}

	if fromAccount.CurrencyCode != toAccount.CurrencyCode {
		tx.Rollback()
		return nil, nil, nil, errors.New("currency mismatch between accounts")
	}

	if fromAccount.Balance < amount {
		tx.Rollback()
		return nil, nil, nil, errors.New("insufficient funds")
	}

	transaction := model.Transaction{
		TransactionReference: generateTransactionReference(),
		TransactionType:      "transfer",
		Amount:               amount,
		CurrencyCode:         fromAccount.CurrencyCode,
		Status:               "completed",
		Description:          description,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := s.TransactionRepo.Create(tx, &transaction); err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	fromAccount.Balance -= amount
	fromAccount.AvailableBalance -= amount
	fromAccount.UpdatedAt = time.Now()

	toAccount.Balance += amount
	toAccount.AvailableBalance += amount
	toAccount.UpdatedAt = time.Now()

	debitEntry := model.LedgerEntry{
		TransactionID:  transaction.ID,
		AccountID:      fromAccount.ID,
		EntryType:      "debit",
		Amount:         amount,
		RunningBalance: fromAccount.Balance,
		CreatedAt:      time.Now(),
	}

	creditEntry := model.LedgerEntry{
		TransactionID:  transaction.ID,
		AccountID:      toAccount.ID,
		EntryType:      "credit",
		Amount:         amount,
		RunningBalance: toAccount.Balance,
		CreatedAt:      time.Now(),
	}

	if err := s.TransactionRepo.CreateLedgerEntry(tx, &debitEntry); err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	if err := s.TransactionRepo.CreateLedgerEntry(tx, &creditEntry); err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	if _, err := s.AccountRepo.Update(tx, fromAccount); err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	if _, err := s.AccountRepo.Update(tx, toAccount); err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	return &transaction, fromAccount, toAccount, nil
}

func generateTransactionReference() string {
	return fmt.Sprintf("TRX-%s", uuid.New().String()[:8])
}
