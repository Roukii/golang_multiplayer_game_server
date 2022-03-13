package world

import (
	"github.com/Roukii/pock_multiplayer/internal/gateway/dao"
	"github.com/Roukii/pock_multiplayer/internal/gateway/entity"
)

type WorldService struct {
	WorldDao *dao.WorldDao
}

// New -.
func NewWorld(r *dao.WorldDao) *WorldService {
	return &WorldService{
		WorldDao: r,
	}
}

func (a *WorldService) GetById(WorldId string) (entity.World, error) {
	world, err := a.WorldDao.GetById(WorldId)
	if err != nil {
		return entity.World{}, err
	}
	return world, err
}

func (a *WorldService) GetUserWorlds(userId string) ([]entity.World, error) {
	worlds, err := a.WorldDao.GetUserWorlds(userId)
	if err != nil {
		return []entity.World{}, err
	}
	return worlds, err
}

func (a *WorldService) GetAvailableWorld() ([]entity.World, error) {
	worlds, err := a.WorldDao.GetAll()
	if err != nil {
		return []entity.World{}, err
	}
	return worlds, err
}
