package service

import (
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/scylladb/gocqlx/v2"
)

type Service struct {
	Logger      logger.Interface
	GameService *game.GameService
}

func New(session *gocqlx.Session, l logger.Interface) *Service {
	return &Service{
		Logger:      l,
		GameService: game.NewGameService("universe_uuid", session),
	}
}
