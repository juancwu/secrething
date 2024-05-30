package entry

import (
	"time"
)

type PersonalBentoEntry struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Content         string    `json:"content"`
	PersonalBentoId string    `json:"personal_bento_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
