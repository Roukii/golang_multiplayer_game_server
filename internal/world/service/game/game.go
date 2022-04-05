package game

import (
	"fmt"

	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/scylladb/gocqlx/v2"
)

type GameService struct {
	Universe        universe.Universe
	ConnectedPlayer []player.Player
	ChunkDao        *dao.ChunkDao
	PlayerDao       *dao.PlayerDao
	WorldDao        *dao.WorldDao
	Logger          *logger.Logger
	WorldGenerators map[string]*procedural_generation.WorldGenerator
}

func NewGameService(universeUUID string, session *gocqlx.Session) *GameService {
	tmp := &GameService{
		Universe: universe.Universe{
			UUID:   universeUUID,
			Worlds: make(map[string]universe.World),
		},
		ChunkDao:        dao.NewChunkDao(session),
		PlayerDao:       dao.NewPlayerDao(session),
		WorldDao:        dao.NewWorldDao(session),
		WorldGenerators: make(map[string]*procedural_generation.WorldGenerator),
	}
	err := tmp.startGame()
	if err != nil {
		fmt.Println("error : ", err)
		return tmp
	}
	return tmp
}

func (g *GameService) startGame() (err error) {
	worlds, err := g.WorldDao.GetAllWorlds()
	for _, world := range worlds {
		g.Universe.Worlds[world.UUID] = world
		generator := procedural_generation.NewWorldGenerator(&world)
		g.WorldGenerators[world.UUID] = &generator
	}
	if err != nil {
		fmt.Println("error : ", err)
		return err
	}

	if len(g.Universe.Worlds) == 0 {
		world, err := g.CreateWorld("tutorial")
		if err != nil {
			fmt.Println("error : ", err)
		}
		fmt.Println("create world with uuid :", world.UUID)
		g.Universe.Worlds[world.UUID] = world
	} else {
		world, err := g.CreateWorld("tutorial")
		if err != nil {
			fmt.Println("error : ", err)
		}
		fmt.Println("world uuid :", world.UUID)
		g.Universe.Worlds = make(map[string]universe.World)
		g.Universe.Worlds[world.UUID] = world
	}
	return err
}
