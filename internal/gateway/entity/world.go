package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type World struct {
	UUID      string `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	UserId    string
	CreatedAt time.Time
	UpdateAt  time.Time
}

func (u *World) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New().String()
	u.CreatedAt = time.Now()
	return nil
}

func (u *World) BeforeSave(db *gorm.DB) (err error) {
	u.UpdateAt = time.Now()
	return nil
}
