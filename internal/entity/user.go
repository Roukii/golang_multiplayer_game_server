package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UUID      string    `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	Username  string    `gorm:"type:varchar(40);unique" json:"username"`
	Password  []byte    `json:"password,omitempty"`
	Worlds    []World   `gorm:"foreignKey:user_id"`
	CreatedAt time.Time `json:"createdAt"`
}


func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  u.UUID = uuid.New().String()
	u.CreatedAt = time.Now()
	return nil
}
