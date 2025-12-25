package database

import (
	"log"
	"paygo/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

// Make It Idempotent (run multiple times safely)

func (d *Database) Seed() error {
	log.Println("Starting database seeding...")

	// Create test users
	users := []model.User{
		{
			ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
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
			ID:           uuid.MustParse("22222222-2222-2222-2222-222222222222"),
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
			ID:               uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
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
			ID:               uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
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

	log.Println("\nDatabase seeding completed successfully!")
	log.Println("\nTest Accounts:")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("Alice's Account: %s (USD 1,000.00)", accounts[0].ID)
	log.Printf("Bob's Account:   %s (USD 500.00)", accounts[1].ID)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("\nYou can now test transfers between these accounts!")

	return nil
}
