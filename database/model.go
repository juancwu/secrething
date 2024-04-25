package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	PublicKey string
	CreatedAt time.Time
	UpdatedAt time.Time
}
