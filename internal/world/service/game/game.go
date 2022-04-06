package game

import (
	"fmt"
	"sync"

	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/scylladb/gocqlx/v2"
)

type GameService struct {
	Universe            universe.Universe
	ConnectedPlayer     map[string]player.Player
	ChunkDao            *dao.ChunkDao
	PlayerDao           *dao.PlayerDao
	WorldDao            *dao.WorldDao
	Logger              *logger.Logger
	PlayerMu            sync.RWMutex
	PlayerActionChannel chan PlayerAction
	PlayerChangeChannel chan PlayerChange
	WorldGenerators     map[string]*procedural_generation.WorldGenerator
}

func NewGameService(universeUUID string, session *gocqlx.Session) *GameService {
	tmp := &GameService{
		Universe: universe.Universe{
			UUID:   universeUUID,
			Worlds: make(map[string]universe.World),
		},
		ChunkDao:            dao.NewChunkDao(session),
		PlayerDao:           dao.NewPlayerDao(session),
		WorldDao:            dao.NewWorldDao(session),
		WorldGenerators:     make(map[string]*procedural_generation.WorldGenerator),
		PlayerChangeChannel: make(chan PlayerChange, 1),
		PlayerActionChannel: make(chan PlayerAction, 1),
		ConnectedPlayer:     make(map[string]player.Player),
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
	}

	go g.watchPlayerActions()
	return err
}

func (g *GameService) watchPlayerActions() {
	for {
		action := <-g.PlayerActionChannel
		g.PlayerMu.Lock()
		action.Perform(g)
		g.PlayerMu.Unlock()
	}
}
