package universe_service

import (
	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	"github.com/scylladb/gocqlx/v2"
)

type UniverseService struct {
	Universe        universe.Universe
	chunkDao        *dao.ChunkDao
	worldDao        *dao.WorldDao
	WorldGenerators map[string]*procedural_generation.WorldGenerator
}

func NewUniverseService(session *gocqlx.Session) *UniverseService {
	u := UniverseService{
		Universe:        universe.Universe{
			Worlds:    make(map[string]universe.World),
		},
		chunkDao:        dao.NewChunkDao(session),
		worldDao:        dao.NewWorldDao(session),
		WorldGenerators: map[string]*procedural_generation.WorldGenerator{},
	}
	return &u
}
