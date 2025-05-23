// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

type User struct {
	UserID       UserID `db:"user_id" json:"user_id"`
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"password_hash"`
	FirstName    string `db:"first_name" json:"first_name"`
	LastName     string `db:"last_name" json:"last_name"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	UpdatedAt    string `db:"updated_at" json:"updated_at"`
}

type Vault struct {
	VaultID      VaultID `db:"vault_id" json:"vault_id"`
	VaultName    string  `db:"vault_name" json:"vault_name"`
	VaultOwnerID UserID  `db:"vault_owner_id" json:"vault_owner_id"`
	CreatedAt    string  `db:"created_at" json:"created_at"`
	UpdatedAt    string  `db:"updated_at" json:"updated_at"`
}
