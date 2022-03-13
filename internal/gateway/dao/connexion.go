package dao

import (
	"github.com/Roukii/pock_multiplayer/internal/gateway/entity"
	"gorm.io/gorm"
)

// ConnexionDao -.
type ConnexionDao struct {
	db *gorm.DB
}

// New -.
func NewConnexionDao(pg *gorm.DB) *ConnexionDao {
	pg.AutoMigrate(entity.Connexion{})
	return &ConnexionDao{pg}
}

func (a ConnexionDao) GetAll() ([]entity.Connexion, error) {
	var conns []entity.Connexion
	result := a.db.Find(conns)
	return conns, result.Error
}

func (a ConnexionDao) SaveOrUpdate(conn *entity.Connexion) error {
	return a.db.Save(conn).Error
}
