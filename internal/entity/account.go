package entity

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Account struct {
	UUID      string    `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt timestamp `json:"createdAt"`
}
