package model

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"userId" gorm:"not null;index"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Balance   float64        `json:"balance" gorm:"default:0.00"`
	Currency  string         `json:"currency" gorm:"size:3;default:USD"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
