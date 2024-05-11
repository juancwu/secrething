package usermodel

import (
	"time"
)

type User struct {
	Id            string
	FirstName     string
	LastName      string
	Password      string // password is always encrypted when fetched from db
	Email         string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
