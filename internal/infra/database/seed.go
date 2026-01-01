package database

import (
	"log"
	"paygo/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

// Make It Idempotent (run multiple times safely)

func timePtr(t time.Time) *time.Time {
	return &t
}

func (d *Database) Seed() error {
	log.Println("Starting database seeding...")

	// Create test users
	users := []model.User{
		{
			ID:           uuid.New(),
			Email:        "alice@example.com",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: "password123"
			FirstName:    "Alice",
			LastName:     "Johnson",
			Verified:     true,
			Status:       "active",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Email:        "bob@example.com",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: "password123"
			FirstName:    "Bob",
			LastName:     "Smith",
			Verified:     true,
			Status:       "active",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, user := range users {
		if err := d.DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
			return err
		}
		log.Printf("Created user: %s", user.Email)
	}

	// Create test accounts (minimum needed for transfer testing)
	accounts := []model.Account{
		{
			ID:               uuid.New(),
			UserID:           users[0].ID,
			AccountNumber:    "ACC-1000001",
			AccountType:      "checking",
			CurrencyCode:     "USD",
			Balance:          1000.00,
			AvailableBalance: 1000.00,
			Status:           "active",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               uuid.New(),
			UserID:           users[1].ID,
			AccountNumber:    "ACC-2000001",
			AccountType:      "checking",
			CurrencyCode:     "USD",
			Balance:          500.00,
			AvailableBalance: 500.00,
			Status:           "active",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}

	for _, account := range accounts {
		if err := d.DB.Create(&account).Error; err != nil {
			log.Printf("Failed to create account %s: %v", account.AccountNumber, err)
			return err
		}
		log.Printf("Created account: %s (Balance: %.2f %s)",
			account.AccountNumber, account.Balance, account.CurrencyCode)
	}

	// Create initial transactions and ledger entries for account balances
	transactions := []model.Transaction{
		{
			ID:                   uuid.New(),
			TransactionReference: "INIT-DEP-001",
			TransactionType:      "deposit",
			Status:               "completed",
			Amount:               1000.00,
			CurrencyCode:         "USD",
			Description:          "Initial deposit",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
		{
			ID:                   uuid.New(),
			TransactionReference: "INIT-DEP-002",
			TransactionType:      "deposit",
			Status:               "completed",
			Amount:               500.00,
			CurrencyCode:         "USD",
			Description:          "Initial deposit",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
	}

	for _, txn := range transactions {
		if err := d.DB.Create(&txn).Error; err != nil {
			log.Printf("Failed to create transaction %s: %v", txn.ID, err)
			return err
		}
		log.Printf("Created initial transaction: %s (%.2f %s)", txn.ID, txn.Amount, txn.CurrencyCode)
	}

	// Create ledger entries for the initial deposits
	ledgerEntries := []model.LedgerEntry{
		{
			ID:             uuid.New(),
			TransactionID:  transactions[0].ID,
			AccountID:      accounts[0].ID,
			EntryType:      "credit",
			Amount:         1000.00,
			RunningBalance: 1000.00,
			CreatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			TransactionID:  transactions[1].ID,
			AccountID:      accounts[1].ID,
			EntryType:      "credit",
			Amount:         500.00,
			RunningBalance: 500.00,
			CreatedAt:      time.Now(),
		},
	}

	for _, entry := range ledgerEntries {
		if err := d.DB.Create(&entry).Error; err != nil {
			log.Printf("Failed to create ledger entry %s: %v", entry.ID, err)
			return err
		}
		log.Printf("Created ledger entry: %s (%s %.2f)",
			entry.AccountID, entry.EntryType, entry.Amount)
	}

	log.Println("\nDatabase seeding completed successfully!")
	log.Println("\nTest Accounts:")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("Alice's Account: %s (USD 1,000.00)", accounts[0].ID)
	log.Printf("Bob's Account:   %s (USD 500.00)", accounts[1].ID)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("\nYou can now test transfers between these accounts!")

	return nil
}
