package models

import "time"

// Challenge represents a randomly generated data that will be used
// for a client to sign it using their private key. Each generated Challenge
// has an expiration time of 30 seconds. This ensures that all challenges
// must be done in a timely matter and that if a signed challenge has been
// stolen, it won't cause harm since there would be a state attached to the challenge.
type Challenge struct {
	id          int64
	State       string `json:"state"`
	HashedValue string `json:"hashed_value"`
	UserId      string
	CreatedAt   time.Time
	ExpiresAt   time.Time `json:"expires_at"`
}
