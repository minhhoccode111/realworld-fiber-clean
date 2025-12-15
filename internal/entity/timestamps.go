package entity

import "time"

// Timestamps tracks lifecycle times for an entity.
type Timestamps struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}
