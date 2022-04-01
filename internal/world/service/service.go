package service

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	entity "github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type Service struct {
	Logger      logger.Interface
	UniverseDao dao.UniverseDao
}

func New(session *gocqlx.Session, l logger.Interface) *Service {
	universeDao := dao.NewUniverseDao(session)
	worlds := []universe.World{}
	worlds = append(worlds, universe.World{
		UUID:      gocql.TimeUUID().String(),
		Name:      "uuid",
		Level:     0,
		Length:    0,
		Width:     0,
		Chunks:    []universe.Chunk{},
		Seed:      "",
		Type:      0,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	})
	if err := universeDao.Insert(&universe.Universe{
		UUID:      gocql.TimeUUID().String(),
		Name:      "Coucou",
		Worlds:    worlds,
		Players:   []entity.Player{},
		CreatedAt: time.Now(),
	}); err != nil {
		l.Debug(err)
		return nil
	}
	universe, err := universeDao.Query()
	if err != nil {
		l.Debug(err)
		return nil
	} else {
		l.Debug("universe", universe)
	}

	return &Service{
		Logger: l,
	}
}
