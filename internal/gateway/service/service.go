package service

import (
	"github.com/Roukii/pock_multiplayer/internal/gateway/dao"
	"github.com/Roukii/pock_multiplayer/internal/gateway/service/user"
	"github.com/Roukii/pock_multiplayer/internal/gateway/service/world"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"gorm.io/gorm"
)

type Service struct {
	UserService          *user.UserService
	WorldService         *world.WorldService
	UserWorldAffdService *world.UserWorldAffService
	Logger               logger.Interface
}

func New(pg *gorm.DB, l logger.Interface) *Service {
	return &Service{
		UserService:          user.New(dao.NewUserDao(pg), dao.NewConnexionDao(pg)),
		WorldService:         world.NewWorld(dao.NewWorldDao(pg)),
		UserWorldAffdService: world.NewUserWorldAff(dao.NewUserWorldAffDao(pg)),
		Logger:               l,
	}
}
