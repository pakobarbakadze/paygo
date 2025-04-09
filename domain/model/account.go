package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID               uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	AccountNumber    string        `gorm:"uniqueIndex;not null" json:"account_number"`
	AccountType      string        `gorm:"not null" json:"account_type"`
	CurrencyCode     string        `gorm:"type:char(3);not null" json:"currency_code"`
	Balance          float64       `gorm:"type:numeric(19,4);not null;default:0" json:"balance"`
	AvailableBalance float64       `gorm:"type:numeric(19,4);not null;default:0" json:"available_balance"`
	Status           string        `gorm:"default:active" json:"status"`
	CreatedAt        time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time     `gorm:"not null" json:"updated_at"`
	User             User          `gorm:"foreignKey:UserID" json:"-"`
	LedgerEntries    []LedgerEntry `gorm:"foreignKey:AccountID" json:"ledger_entries,omitempty"`
}
