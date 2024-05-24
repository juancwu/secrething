package bentomodel

import "github.com/juancwu/konbini/server/database"

func PersonalBentoExistsWithName(userId, name string) (bool, error) {
	var exists bool
	err := database.DB().QueryRow(
		"SELECT EXISTS ( SELECT 1 FROM personal_bentos WHERE owner_id = $1 AND name = $2)",
		userId,
		name,
	).Scan(&exists)
	return exists, err
}

func NewPersonalBento(ownerId, name, pubKey, content string) (string, error) {
	var id string
	row := database.DB().QueryRow("INSERT INTO personal_bentos (owner_id, name, pub_key, content) VALUES ($1, $2, $3, $4) RETURNING id;", ownerId, name, pubKey, content)
	err := row.Scan(
		&id,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}
