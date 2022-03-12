package dao

import (
	"github.com/Roukii/pock_multiplayer/internal/entity"
	"gorm.io/gorm"
)

// UserDao -.
type UserDao struct {
	db *gorm.DB
}

// New -.
func NewUserDao(pg *gorm.DB) *UserDao {
	pg.AutoMigrate(entity.User{})
	return &UserDao{pg}
}

func (a UserDao) GetById(userId string) (entity.User, error) {
	var user entity.User
	result := a.db.First(&user, userId)
	return user, result.Error
}

func (a UserDao) GetAll() ([]entity.User, error) {
	var users []entity.User
	result := a.db.Find(users)
	return users, result.Error
}

func (a UserDao) SaveOrUpdate(user *entity.User) error {
	return a.db.Save(user).Error
}
