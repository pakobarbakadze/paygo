package model

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID                   uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionReference string        `gorm:"uniqueIndex;not null" json:"transaction_reference"`
	TransactionType      string        `gorm:"not null" json:"transaction_type"`
	Amount               float64       `gorm:"type:numeric(19,4);not null" json:"amount"`
	CurrencyCode         string        `gorm:"type:char(3);not null" json:"currency_code"`
	Status               string        `gorm:"default:pending" json:"status"`
	Description          string        `json:"description"`
	CreatedAt            time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt            time.Time     `gorm:"not null" json:"updated_at"`
	LedgerEntries        []LedgerEntry `gorm:"foreignKey:TransactionID" json:"ledger_entries,omitempty"`
}
