package dao

import (
	"github.com/Roukii/pock_multiplayer/internal/gateway/entity"
	"gorm.io/gorm"
)

// WorlDao -.
type UserWorldAffDao struct {
	*gorm.DB
}

// New -.
func NewUserWorldAffDao(pg *gorm.DB) *UserWorldAffDao {
	pg.AutoMigrate(entity.UserWorldAff{})
	return &UserWorldAffDao{pg}
}

func (a *UserWorldAffDao) Create(aff *entity.UserWorldAff) error {
	return a.DB.Create(aff).Error
}

func (a *UserWorldAffDao) Delete(aff *entity.UserWorldAff) error {
	return a.DB.Delete(aff).Error
}
