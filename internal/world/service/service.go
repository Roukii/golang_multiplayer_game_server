package service

import (
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/gocql/gocql"
)

type Service struct {
	Logger logger.Interface
}

func New(session *gocql.Session, l logger.Interface) *Service {
	return &Service{
		Logger: l,
	}
}
