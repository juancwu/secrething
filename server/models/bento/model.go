package bentomodel

import "time"

type PersonalBento struct {
	Id        string    `json:"id"`
	OwnerId   string    `json:"owner_id"`
	Name      string    `json:"name"`
	Content   []byte    `json:"content"`
	PubKey    []byte    `json:"pub_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
