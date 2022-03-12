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
	// TODO join table user with user_world_aff
	// db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
	return []entity.World{}, nil
}
