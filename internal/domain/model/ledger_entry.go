package model

import (
	"time"

	"github.com/google/uuid"
)

type LedgerEntry struct {
	ID             uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionID  uuid.UUID   `gorm:"type:uuid;not null" json:"transaction_id"`
	AccountID      uuid.UUID   `gorm:"type:uuid;not null" json:"account_id"`
	EntryType      string      `gorm:"not null" json:"entry_type"` // "debit" or "credit"
	Amount         float64     `gorm:"type:numeric(19,4);not null" json:"amount"`
	RunningBalance float64     `gorm:"type:numeric(19,4);not null" json:"running_balance"`
	CreatedAt      time.Time   `gorm:"not null" json:"created_at"`
	Transaction    Transaction `gorm:"foreignKey:TransactionID" json:"-"`
	Account        Account     `gorm:"foreignKey:AccountID" json:"-"`
}
