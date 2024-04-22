package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Conn *gorm.DB
}

func New() *DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	return &DB{Conn: db}
}

func (db *DB) Migrate() {
	db.Conn.AutoMigrate(&User{})
}

func (db *DB) CreateUser(name, email, publicKey string) {
	result := db.Conn.Create(&User{Name: name, Email: email, PublicKey: publicKey})
	if result.Error != nil {
		log.Fatalf("Error creating user: %v\n", result.Error)
	}

	fmt.Printf("Rows affected: %d\n", result.RowsAffected)
}
