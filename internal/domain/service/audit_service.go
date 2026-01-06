package service

import (
	"fmt"
	"paygo/internal/domain/model"
	"paygo/internal/domain/repository"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AuditStatus string

const (
	AuditStatusValid      AuditStatus = "VALID"
	AuditStatusFraudulent AuditStatus = "FRAUDULENT"
	AuditStatusIncomplete AuditStatus = "INCOMPLETE"
)

type FraudType string

const (
	FraudTypeBalanceMismatch FraudType = "BALANCE_MISMATCH"
)

type AuditResult struct {
	AccountID          uuid.UUID   `json:"account_id"`
	AccountNumber      string      `json:"account_number"`
	Status             AuditStatus `json:"status"`
	ExpectedBalance    float64     `json:"expected_balance"`
	ActualBalance      float64     `json:"actual_balance"`
	BalanceDiscrepancy float64     `json:"balance_discrepancy"`
	FraudTypes         []FraudType `json:"fraud_types,omitempty"`
	Details            []string    `json:"details,omitempty"`
	LedgerEntriesCount int         `json:"ledger_entries_count"`
	AuditedAt          time.Time   `json:"audited_at"`
}

type AuditService struct {
	accountRepo *repository.AccountRepository
}

func NewAuditService(
	accountRepo *repository.AccountRepository,
) *AuditService {
	return &AuditService{
		accountRepo: accountRepo,
	}
}

func (s *AuditService) AuditAccounts(accountIDs []uuid.UUID) []AuditResult {
	resultChan := make(chan AuditResult, len(accountIDs))
	var wg sync.WaitGroup

	for _, id := range accountIDs {
		wg.Add(1)
		go func(accountID uuid.UUID) {
			defer wg.Done()
			resultChan <- s.AuditAccount(accountID)
		}(id)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	results := make([]AuditResult, 0, len(accountIDs))
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

func (s *AuditService) AuditAccount(accountID uuid.UUID) AuditResult {
	result := AuditResult{
		AccountID:  accountID,
		Status:     AuditStatusValid,
		FraudTypes: []FraudType{},
		Details:    []string{},
		AuditedAt:  time.Now(),
	}

	account, err := s.accountRepo.FindByID(accountID, false)
	if err != nil {
		result.Status = AuditStatusIncomplete
		result.Details = append(result.Details, fmt.Sprintf("Failed to fetch account: %v", err))
		return result
	}

	result.AccountNumber = account.AccountNumber
	result.ActualBalance = account.Balance
	result.LedgerEntriesCount = len(account.LedgerEntries)

	expectedBalance := s.calculateExpectedBalanceFromLedger(account)
	result.ExpectedBalance = expectedBalance
	result.BalanceDiscrepancy = account.Balance - expectedBalance

	if account.Balance != expectedBalance {
		result.Status = AuditStatusFraudulent
		result.FraudTypes = append(result.FraudTypes, FraudTypeBalanceMismatch)
		result.Details = append(result.Details, fmt.Sprintf(
			"Balance mismatch detected: actual=%.4f, expected=%.4f, discrepancy=%.4f",
			account.Balance, expectedBalance, result.BalanceDiscrepancy,
		))
	}

	return result
}

func (s *AuditService) calculateExpectedBalanceFromLedger(account *model.Account) float64 {
	balance := 0.0

	for _, entry := range account.LedgerEntries {
		if entry.EntryType == "credit" {
			balance += entry.Amount
		} else if entry.EntryType == "debit" {
			balance -= entry.Amount
		}
	}

	return balance
}
