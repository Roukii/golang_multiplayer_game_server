package entity

import "time"

// Translation -.
type World struct {
	UUID      string    `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	CreatedAt time.Time `json:"createdAt"`
}
