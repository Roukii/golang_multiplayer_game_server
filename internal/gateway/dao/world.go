package dao

import (
	"github.com/Roukii/pock_multiplayer/internal/gateway/entity"
	"gorm.io/gorm"
)

// WorlDao -.
type WorldDao struct {
	*gorm.DB
}

// New -.
func NewWorldDao(pg *gorm.DB) *WorldDao {
	pg.AutoMigrate(entity.World{})
	pg.Create(&entity.World{UUID: "898eddb5-03b3-494d-8adb-f405af75d323", Name: "Valhalla", PlayerCount: 0, MaxPlayer: 10, IsAcceptingPlayer: true})
	pg.Create(&entity.World{UUID: "dc2a56c3-9468-4ce6-bd9f-8a0b80958d4c", Name: "England", PlayerCount: 0, MaxPlayer: 10, IsAcceptingPlayer: true})
	pg.Create(&entity.World{UUID: "1646126a-b454-4312-a169-d54af47bf3fc", Name: "Neerlander", PlayerCount: 0, MaxPlayer: 10, IsAcceptingPlayer: true})
	return &WorldDao{pg}
}

func (a *WorldDao) GetById(worldId string) (entity.World, error) {
	var world entity.World
	result := a.DB.First(&world, worldId)
	return world, result.Error
}

func (a *WorldDao) GetAll() ([]entity.World, error) {
	var worlds []entity.World
	result := a.DB.Find(&worlds)
	return worlds, result.Error
}

func (a *WorldDao) SaveOrUpdate(world *entity.World) error {
	return a.DB.Save(world).Error
}

func (a *WorldDao) GetUserWorlds(userId string) ([]entity.World, error) {
	var worlds []entity.World
	result := a.DB.Model(&entity.World{}).Joins("join user_world_affs uw on uw.world_uuid = worlds.uuid and uw.user_uuid = ?", userId).Scan(&worlds)
	return worlds, result.Error
}
