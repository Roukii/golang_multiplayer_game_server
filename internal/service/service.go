package service

import (
	"github.com/Roukii/pock_multiplayer/internal/dao"
	"github.com/Roukii/pock_multiplayer/internal/service/user"
	"github.com/Roukii/pock_multiplayer/internal/service/world"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"gorm.io/gorm"
)

type Service struct {
	userService  *user.UserService
	worldService *world.WorldService
}

func New(pg *gorm.DB, l logger.Interface) *Service {
	return &Service{
		userService:  user.New(dao.NewUserDao(pg)),
		worldService: world.New(dao.NewWorldDao(pg)),
	}
}
