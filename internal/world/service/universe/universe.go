package universe_service

import (
	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	"github.com/scylladb/gocqlx/v2"
)

type UniverseService struct {
	Universe        universe.Universe
	ChunkDao        *dao.ChunkDao
	WorldGenerators map[string]*procedural_generation.WorldGenerator
	WorldDao        *dao.WorldDao
}

func NewUniverseService(session *gocqlx.Session) *UniverseService {
	u := UniverseService{
		Universe:        universe.Universe{},
		ChunkDao:        dao.NewChunkDao(session),
		WorldDao:        dao.NewWorldDao(session),
		WorldGenerators: map[string]*procedural_generation.WorldGenerator{},
	}
	return &u
}
