package bentomodel

import "time"

type PersonalBento struct {
	Id        string    `json:"id"`
	OwnerId   string    `json:"owner_id"`
	Name      string    `json:"name"`
	KeyVals   []string  `json:"keyvals"`
	PubKey    string    `json:"pub_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
