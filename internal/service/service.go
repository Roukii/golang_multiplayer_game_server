package service

import (
	"github.com/Roukii/pock_multiplayer/internal/dao"
	"github.com/Roukii/pock_multiplayer/internal/service/user"
	"github.com/Roukii/pock_multiplayer/internal/service/world"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"gorm.io/gorm"
)

type Service struct {
	UserService  *user.UserService
	WorldService *world.WorldService
	Logger			 logger.Interface
}

func New(pg *gorm.DB, l logger.Interface) *Service {
	return &Service{
		UserService:  user.New(dao.NewUserDao(pg)),
		WorldService: world.New(dao.NewWorldDao(pg)),
		Logger: l,
	}
}
