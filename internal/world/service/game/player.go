package game

import (
	"fmt"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/gocql/gocql"
)

func (g *GameService) CreatePlayer(userUuid string, p *player.Player) (err error) {
	p.UUID = gocql.TimeUUID().String()
	p.SpawnPoint, err = g.GenerateSpawnPoint(p)
	if err != nil {
		return err
	}
	fmt.Println(p.SpawnPoint)
	err = g.PlayerDao.Insert(userUuid, p)
	if err != nil {
		return err
	}
	g.ConnectedPlayer = append(g.ConnectedPlayer, *p)
	return nil
}

func (g *GameService) ConnectPlayer(userUUID string, playerUUID string) (*player.Player, error) {
	player, err := g.PlayerDao.GetPlayerFromUUID(userUUID, playerUUID)
	if err != nil {
		return nil, err
	}
	g.ConnectedPlayer = append(g.ConnectedPlayer, *player)
	return player, err
}

func (g *GameService) DisconnectPlayer(p *player.Player) (bool, error) {
	err := g.PlayerDao.Update(p)
	if err != nil {
		return false, nil
	}
	for index, v := range g.ConnectedPlayer {
		if v.UUID == p.UUID {
			g.ConnectedPlayer = append(g.ConnectedPlayer[:index], g.ConnectedPlayer[index+1:]...)
			return true, err
		}
	}
	return false, err
}

// TODO choose a planet at random
func (g *GameService) GenerateSpawnPoint(p *player.Player) (player.SpawnPoint, error) {
	var worldUUID string
	for _, world := range g.Universe.Worlds {
		worldUUID = world.UUID
		break
	}
	spawnPoint := player.SpawnPoint{
		WorldUUID: worldUUID,
		Coordinate: entity.Position{
			Position: entity.Vector3f{
				X: 10,
				Y: 10,
				Z: 10,
			},
		},
		UpdatedAt: time.Time{},
	}
	return spawnPoint, nil
}
