// This file contains all things related to how the database store users
package store

import "time"

// User is someone who uses konbini services. This structure is a representation of
// a complete user fetched from the database.
type User struct {
	Id            string
	FirstName     string
	LastName      string
	Email         string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
