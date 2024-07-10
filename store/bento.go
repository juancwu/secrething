package store

import "time"

// Bento represents a bento object that has a slice with all the entries as well.
type Bento struct {
	Id        string
	Name      string
	OwnerId   string
	PubKey    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBento will create and save a new bento into the database with the given information.
// This method will return an error if there is another bento with the same name from the same user.
// All bentos belonging to one user should have unique names.
func NewBento(name, ownerId, pubKey string) (*Bento, error) {
	row := db.QueryRow(
		"INSERT INTO bentos (name, owner_id, pub_key) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at;",
		name,
		ownerId,
		pubKey,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	bento := Bento{
		Name:    name,
		OwnerId: ownerId,
		PubKey:  pubKey,
	}
	err = row.Scan(
		&bento.Id,
		&bento.CreatedAt,
		&bento.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &bento, nil
}
