package entity

import "time"

type World struct {
	UUID      string    `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	UserId    string    `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

