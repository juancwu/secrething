package bentomodel

import "time"

type PersonalBento struct {
	Id        string
	OwnerId   string
	Name      string
	Content   string
	PubKey    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
