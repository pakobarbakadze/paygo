package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	FirstName  string         `json:"firstName" gorm:"size:100;not null"`
	LastName   string         `json:"lastName" gorm:"size:100;not null"`
	Email      string         `json:"email" gorm:"size:255;not null;uniqueIndex"`
	Password   string         `json:"-" gorm:"size:255;not null"`
	Balance    float64        `json:"balance" gorm:"default:0.00"`
	IsVerified bool           `json:"isVerified" gorm:"default:false"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
