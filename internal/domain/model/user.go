package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Verified     bool      `gorm:"default:false" json:"verified"`
	Status       string    `gorm:"default:active" json:"status"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at"`
	Accounts     []Account `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
}
