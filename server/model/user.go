package model

type UserModel struct {
	ID           int
	FirstName    *string
	LastName     *string
	Email        string
	PemPublicKey string
}
