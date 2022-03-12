package entity

import (
	"time"
)

type User struct {
	UUID      string    `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	Username  string    `gorm:"type:varchar(40);unique" json:"username"`
	Password  []byte    `json:"password,omitempty"`
	Worlds    []World   `gorm:"foreignKey:user_id"`
	CreatedAt time.Time `json:"createdAt"`
}
