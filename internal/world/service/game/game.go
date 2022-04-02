package game

import (
	"fmt"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type GameService struct {
	Universe        universe.Universe
	ConnectedPlayer []player.Player
	ChunkDao        *dao.ChunkDao
	PlayerDao       *dao.PlayerDao
	WorldDao        *dao.WorldDao
	logger          *logger.Logger
}

func NewGameService(universeUUID string, session *gocqlx.Session) *GameService {
	tmp := &GameService{
		Universe: universe.Universe{
			UUID: universeUUID,
		},
		ChunkDao:  dao.NewChunkDao(session),
		PlayerDao: dao.NewPlayerDao(session),
		WorldDao:  dao.NewWorldDao(session),
	}
	err := tmp.StartGame()
	if err != nil {
		fmt.Println("error : ", err)
		return tmp
	}
	return tmp
}

func (g *GameService) StartGame() (err error) {
	g.Universe.Worlds, err = g.WorldDao.GetAllWorlds()
	if err != nil {
		fmt.Println("error : ", err)
		return err
	}
	if len(g.Universe.Worlds) == 0 {
		err = g.WorldDao.Insert(&universe.World{
			UUID:      gocql.TimeUUID().String(),
			Name:      "Test",
			MaxPlayer: 20,
			Level:     1,
			Length:    10,
			Width:     10,
			Seed:      "weshletest",
			Type:      0,
			CreatedAt: time.Now(),
			UpdateAt:  time.Now(),
		})
		if err != nil {
			fmt.Println("error : ", err)
			return err
		}
	} else {
		fmt.Println("worlds uuid : ", g.Universe.Worlds[0].UUID)
	}
	return err
}

func (g *GameService) CreatePlayer(userUuid string, p *player.Player) error {
	err := g.PlayerDao.Insert(userUuid, p)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameService) ConnectPlayer(playerUUID player.Player) error {
	return nil
}

func (g *GameService) SpawnPlayer(p *player.Player) (spawnPoint player.SpawnPoint, err error) {
	return spawnPoint, err
}
