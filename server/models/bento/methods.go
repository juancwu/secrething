package bentomodel

import (
	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/utils"
)

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

func GetPersonalBento(uuid string) (*PersonalBento, error) {
	bento := PersonalBento{}
	err := database.DB().QueryRow("SELECT id, name, owner_id, content, pub_key, created_at, updated_at FROM personal_bentos WHERE id = $1;", uuid).Scan(
		&bento.Id,
		&bento.Name,
		&bento.OwnerId,
		&bento.Content,
		&bento.PubKey,
		&bento.CreatedAt,
		&bento.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &bento, nil
}

// IMPORTANT: Only fills out the id, name, created_at and updated_at fields
func ListPersonalBentos(uid string) ([]PersonalBento, error) {
	var bentos []PersonalBento
	rows, err := database.DB().Query("SELECT id, name, created_at, updated_at FROM personal_bentos WHERE owner_id = $1;", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		bento := PersonalBento{}
		err = rows.Scan(&bento.Id, &bento.Name, &bento.CreatedAt, &bento.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bentos = append(bentos, bento)
	}

	return bentos, nil
}

func DeletePersonalBento(uid, bid string) (bool, error) {
	result, err := database.DB().Exec("DELETE FROM personal_bentos WHERE owner_id = $1 AND id = $2;", uid, bid)
	if err != nil {
		return false, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	// shouldn't really happen because the id and owner id should be unique in the database
	if n > 1 {
		utils.Logger().Warn("MORE THAN 1 PERSONAL WAS DELETED!!!!", "uid", uid, "bid", bid)
	}

	return n > 0, nil
}
