package world

import (
	"github.com/Roukii/pock_multiplayer/internal/gateway/dao"
	"github.com/Roukii/pock_multiplayer/internal/gateway/entity"
)

type UserWorldAffService struct {
	UserWorldAffDao *dao.UserWorldAffDao
}

// New -.
func NewUserWorldAff(r *dao.UserWorldAffDao) *UserWorldAffService {
	return &UserWorldAffService{
		UserWorldAffDao: r,
	}
}

func (a *UserWorldAffService) JoinWorld(userId string, worldId string) error {
	aff := entity.UserWorldAff{
		WorldUUID: worldId,
		UserUUID:  userId,
	}
	return a.UserWorldAffDao.Create(&aff)
}

func (a *UserWorldAffService) DeleteWorld(userId string, worldId string) error {
	aff := entity.UserWorldAff{
		WorldUUID: worldId,
		UserUUID:  userId,
	}
	return a.UserWorldAffDao.Delete(&aff)
}
