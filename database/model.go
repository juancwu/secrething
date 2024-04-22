package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	PublicKey string
}
