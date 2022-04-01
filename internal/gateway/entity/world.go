package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type World struct {
	UUID              string `gorm:"column:uuid;primary_key;type:varchar(64)" json:"uuid"`
	Name              string
	PlayerCount       int       `json:"player_count"`
	MaxPlayer         int       `json:"max_player"`
	IsAcceptingPlayer bool      `json:"is_accepting_player"`
	CreatedAt         time.Time `json:"created_at"`
	UpdateAt          time.Time `json:"updatedAt"`
}

func (u *World) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UUID == "" {
		u.UUID = uuid.New().String()
		u.CreatedAt = time.Now()
		u.UpdateAt = time.Now()
	}
	return nil
}

func (u *World) BeforeSave(db *gorm.DB) (err error) {
	u.UpdateAt = time.Now()
	return nil
}
