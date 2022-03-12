package service

import (
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"gorm.io/gorm"
)

type Service struct {
	Logger logger.Interface
}

func New(pg *gorm.DB, l logger.Interface) *Service {
	return &Service{
		Logger: l,
	}
}
