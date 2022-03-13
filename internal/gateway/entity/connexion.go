package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Connexion struct {
	UUID      string    `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	UserId    string    `json:"user_id"`
	Ip        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}

func (u *Connexion) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New().String()
	u.CreatedAt = time.Now()
	return nil
}

func (u *Connexion) BeforeSave(db *gorm.DB) (err error) {
	u.UpdateAt = time.Now()
	return nil
}
