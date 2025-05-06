package service

import (
	"errors"
	"fmt"
	"paygo/internal/domain/model"
	"paygo/internal/infra/database"
	"time"

	"github.com/google/uuid"
)

type accountRepository interface {
	FindByID(tx database.Transaction, id uuid.UUID, forUpdate bool) (*model.Account, error)
	Update(tx database.Transaction, account *model.Account) (*model.Account, error)
}

type transactionRepository interface {
	Create(tx database.Transaction, transaction *model.Transaction) error
	CreateLedgerEntry(tx database.Transaction, entry *model.LedgerEntry) error
}

type TransferService struct {
	DB              database.DBManager
	AccountRepo     accountRepository
	TransactionRepo transactionRepository
}

func NewTransferService(
	db database.DBManager,
	accountRepo accountRepository,
	transactionRepo transactionRepository,
) *TransferService {
	return &TransferService{
		DB:              db,
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
	}
}

func (s *TransferService) TransferMoney(fromAccountID, toAccountID uuid.UUID, amount float64, description string) (*model.Transaction, *model.Account, *model.Account, error) {
	var fromAccount, toAccount *model.Account
	var transaction model.Transaction

	err := s.DB.WithTransaction(func(tx database.Transaction) error {
		var err error

		if fromAccount, toAccount, err = s.validateAccounts(tx, fromAccountID, toAccountID, amount); err != nil {
			return err
		}

		transaction = s.createTransaction(fromAccount.CurrencyCode, amount, description)

		if err := s.TransactionRepo.Create(tx, &transaction); err != nil {
			return err
		}

		s.updateAccountBalances(fromAccount, toAccount, amount)

		if err := s.createLedgerEntries(tx, &transaction, fromAccount, toAccount, amount); err != nil {
			return err
		}

		if err := s.updateAccounts(tx, fromAccount, toAccount); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return &transaction, fromAccount, toAccount, nil
}

func (s *TransferService) validateAccounts(tx database.Transaction, fromAccountID, toAccountID uuid.UUID, amount float64) (*model.Account, *model.Account, error) {
	fromAccount, err := s.AccountRepo.FindByID(tx, fromAccountID, true)
	if err != nil {
		return nil, nil, err
	}

	toAccount, err := s.AccountRepo.FindByID(tx, toAccountID, true)
	if err != nil {
		return nil, nil, err
	}

	if fromAccount.Status != "active" {
		return nil, nil, errors.New("source account is not active")
	}

	if toAccount.Status != "active" {
		return nil, nil, errors.New("destination account is not active")
	}

	if fromAccount.CurrencyCode != toAccount.CurrencyCode {
		return nil, nil, errors.New("currency mismatch between accounts")
	}

	if fromAccount.Balance < amount {
		return nil, nil, errors.New("insufficient funds")
	}

	return fromAccount, toAccount, nil
}

func (s *TransferService) createTransaction(currencyCode string, amount float64, description string) model.Transaction {
	return model.Transaction{
		TransactionReference: generateTransactionReference(),
		TransactionType:      "transfer",
		Amount:               amount,
		CurrencyCode:         currencyCode,
		Status:               "completed",
		Description:          description,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

func (s *TransferService) updateAccountBalances(fromAccount, toAccount *model.Account, amount float64) {
	fromAccount.Balance -= amount
	fromAccount.AvailableBalance -= amount
	fromAccount.UpdatedAt = time.Now()

	toAccount.Balance += amount
	toAccount.AvailableBalance += amount
	toAccount.UpdatedAt = time.Now()
}

func (s *TransferService) createLedgerEntries(tx database.Transaction, transaction *model.Transaction, fromAccount, toAccount *model.Account, amount float64) error {
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
		return err
	}

	if err := s.TransactionRepo.CreateLedgerEntry(tx, &creditEntry); err != nil {
		return err
	}

	return nil
}

func (s *TransferService) updateAccounts(tx database.Transaction, fromAccount, toAccount *model.Account) error {
	if _, err := s.AccountRepo.Update(tx, fromAccount); err != nil {
		return err
	}

	if _, err := s.AccountRepo.Update(tx, toAccount); err != nil {
		return err
	}

	return nil
}

func generateTransactionReference() string {
	return fmt.Sprintf("TRX-%s", uuid.New().String()[:8])
}
