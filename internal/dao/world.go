package dao

import (
	"github.com/Roukii/pock_multiplayer/internal/entity"
	"gorm.io/gorm"
)

// WorlDao -.
type WorldDao struct {
	*gorm.DB
}

// New -.
func NewWorldDao(pg *gorm.DB) *WorldDao {
	pg.AutoMigrate(entity.User{})
	return &WorldDao{pg}
}

func (a *WorldDao) GetById(worldId string) (entity.World, error) {
	var world entity.World
	result := a.DB.First(&world, worldId)
	return world, result.Error
}

func (a *WorldDao) GetAll() ([]entity.World, error) {
	var worlds []entity.World
	result := a.DB.Find(worlds)
	return worlds, result.Error
}

func (a *WorldDao) SaveOrUpdate(world *entity.World) error {
	return a.DB.Save(world).Error
}

func (a *WorldDao) getUserWorlds(userId string) ([]entity.World, error) {
	// TODO

}
