package model

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Balance   float64   `gorm:"type:numeric(19,4);not null;default:0" json:"balance"`
	Currency  string    `gorm:"type:char(3);not null;default:USD" json:"currency"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
