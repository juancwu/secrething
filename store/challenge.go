package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Challenge struct {
	Id        string
	Challenge string
	BentoId   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewChallenge will create and store a new challenge in the database which has 1 minute lifespan.
// The function returns the challenge id, challange, error.
func NewChallenge(bentoId string) (string, string, error) {
	challenge, err := getNewChallenge()
	if err != nil {
		return "", "", err
	}

	var challengeId string
	exp := time.Now().Add(time.Minute)
	fmt.Println(challenge, bentoId, exp)
	row := db.QueryRow(
		"INSERT INTO challenges (challenge, bento_id, expires_at) VALUES ($1, $2, $3) RETURNING id;",
		challenge,
		bentoId,
		exp,
	)
	if err := row.Err(); err != nil {
		return "", "", err
	}
	err = row.Scan(&challengeId)
	if err != nil {
		return "", "", err
	}

	return challengeId, challenge, nil
}

// GetChallenge retrieves an store challenge in the database using the provided challenge id.
func GetChallenge(challengeId string) (*Challenge, error) {
	challenge := Challenge{}
	row := db.QueryRow(
		"SELECT id, challenge, bento_id, created_at, expires_at FROM challenges WHERE id = $1;",
		challengeId,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	err = row.Scan(
		&challenge.Id,
		&challenge.BentoId,
		&challenge.CreatedAt,
		&challenge.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

// DeleteChallenge deletes a challenge with the given id.
func DeleteChallenge(challengeId string) error {
	_, err := db.Exec(
		"DELETE FROM challenges WHERE id = $1",
		challengeId,
	)
	return err
}

// getNewChallenge generates a new random challenge that has an expiration time and a state attached to it.
// When a client wants to modify a bento, it will need it.
func getNewChallenge() (string, error) {
	challengeBytes := make([]byte, 32)
	_, err := rand.Read(challengeBytes)
	if err != nil {
		return "", err
	}
	_, err = rand.Read(challengeBytes)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(challengeBytes)
	hashedBytes := hash.Sum(nil)
	challenge := hex.EncodeToString(hashedBytes)

	return challenge, nil
}
