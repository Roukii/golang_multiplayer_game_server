package game

import (
	"log"

	player_service "github.com/Roukii/pock_multiplayer/internal/world/service/player"
	universe_service "github.com/Roukii/pock_multiplayer/internal/world/service/universe"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/scylladb/gocqlx/v2"
)

type GameService struct {
	PlayerService       *player_service.PlayerService
	UniverseService     *universe_service.UniverseService
	Logger              *logger.Logger
	PlayerActionChannel chan PlayerAction
	PlayerChangeChannel chan PlayerChange
}

func NewGameService(universeUUID string, session *gocqlx.Session) *GameService {
	tmp := &GameService{
		PlayerService:       player_service.NewPlayerService(session),
		UniverseService:     universe_service.NewUniverseService(session),
		PlayerChangeChannel: make(chan PlayerChange, 1),
		PlayerActionChannel: make(chan PlayerAction, 1),
	}
	err := tmp.startGame()
	if err != nil {
		log.Println("error : ", err)
		return tmp
	}
	return tmp
}

func (g *GameService) startGame() (err error) {
	err = g.UniverseService.LoadWorlds()
	if err != nil {
		log.Println("error : ", err)
		return err
	}

	go g.watchPlayerActions()
	return err
}

func (g *GameService) watchPlayerActions() {
	for {
		action := <-g.PlayerActionChannel
		action.Perform(g)
	}
}
