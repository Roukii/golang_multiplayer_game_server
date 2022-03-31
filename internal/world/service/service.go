package service

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	entity "github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/gocql/gocql"
)

type Service struct {
	Logger      logger.Interface
	UniverseDao dao.UniverseDao
}

func New(session *gocql.Session, l logger.Interface) *Service {
	universeDao := dao.NewUniverseDao(session)
	if err := universeDao.Insert(&universe.Universe{
		UUID:      gocql.TimeUUID().String(),
		Name:      "Coucou",
		Worlds:    []universe.World{},
		Players:   []entity.Player{},
		CreatedAt: time.Time{},
	}); err != nil {
		l.Debug(err)
	}
	return &Service{
		Logger: l,
	}
}
