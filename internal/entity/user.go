package entity

import (
	"time"
)

type User struct {
	UUID      string    `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}
