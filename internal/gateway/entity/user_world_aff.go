package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserWorldAff struct {
	WorldUUID string    `gorm:"primaryKey" json:"world_uuid"`
	UserUUID  string    `gorm:"primaryKey" json:"user_uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updatedAt"`
}

func (u *UserWorldAff) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return nil
}

func (u *UserWorldAff) BeforeSave(db *gorm.DB) (err error) {
	u.UpdateAt = time.Now()
	return nil
}
