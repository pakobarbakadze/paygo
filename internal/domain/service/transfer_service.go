package service

import (
	"errors"
	"fmt"
	"paygo/internal/domain/model"
	"paygo/internal/domain/repository"
	"paygo/internal/infra/database"
	"time"

	"github.com/google/uuid"
)

// TODO: Which is better to use concreted types or interfaces for repos and db manager?
type TransferService struct {
	DB              database.DBManager
	AccountRepo     *repository.AccountRepository
	TransactionRepo *repository.TransactionRepository
}

func NewTransferService(
	db database.DBManager,
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
	var transaction model.Transaction

	// TODO: context is not passed.
	err := s.DB.WithTransaction(func(tx database.DB) error {
		var err error

		txAccountRepo := s.AccountRepo.WithTx(tx)
		txTransactionRepo := s.TransactionRepo.WithTx(tx)

		if fromAccount, toAccount, err = s.validateAccounts(txAccountRepo, fromAccountID, toAccountID, amount); err != nil {
			return err
		}

		transaction = s.createTransaction(fromAccount.CurrencyCode, amount, description)

		if err := txTransactionRepo.Create(&transaction); err != nil {
			return err
		}

		s.updateAccountBalances(fromAccount, toAccount, amount)

		if err := s.createLedgerEntries(txTransactionRepo, &transaction, fromAccount, toAccount, amount); err != nil {
			return err
		}

		if err := s.updateAccounts(txAccountRepo, fromAccount, toAccount); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return &transaction, fromAccount, toAccount, nil
}

func (s *TransferService) validateAccounts(repo *repository.AccountRepository, fromAccountID, toAccountID uuid.UUID, amount float64) (*model.Account, *model.Account, error) {
	fromAccount, err := repo.FindByID(fromAccountID, true)
	if err != nil {
		return nil, nil, err
	}

	toAccount, err := repo.FindByID(toAccountID, true)
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

func (s *TransferService) createLedgerEntries(repo *repository.TransactionRepository, transaction *model.Transaction, fromAccount, toAccount *model.Account, amount float64) error {
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

	if err := repo.CreateLedgerEntry(&debitEntry); err != nil {
		return err
	}

	if err := repo.CreateLedgerEntry(&creditEntry); err != nil {
		return err
	}

	return nil
}

func (s *TransferService) updateAccounts(repo *repository.AccountRepository, fromAccount, toAccount *model.Account) error {
	if _, err := repo.Update(fromAccount); err != nil {
		return err
	}

	if _, err := repo.Update(toAccount); err != nil {
		return err
	}

	return nil
}

func generateTransactionReference() string {
	return fmt.Sprintf("TRX-%s", uuid.New().String()[:8])
}
