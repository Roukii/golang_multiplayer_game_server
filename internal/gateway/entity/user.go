package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UUID      string         `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	Username  string         `gorm:"type:varchar(40);unique" json:"username"`
	Password  []byte         `json:"password,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdateAt  time.Time      `json:"updateAt"`
	Worlds    []UserWorldAff
	Connexion []Connexion
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New().String()
	u.CreatedAt = time.Now()
	return nil
}

func (u *User) BeforeSave(db *gorm.DB) (err error) {
	u.UpdateAt = time.Now()
	return nil
}
