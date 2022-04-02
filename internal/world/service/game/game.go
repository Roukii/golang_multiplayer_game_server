package game

import (
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
)

type game struct {
	Universe        universe.Universe
	ConnectedPlayer []player.Player
}

type Game interface {
	StartGame() error
	AddPlayer(playerUUID player.Player) error
}

func NewGame(universeUUID string) Game {
	return &game{
		Universe: universe.Universe{
			UUID: universeUUID,
		},
	}
}

func (g *game) StartGame() error {
	return nil
}

func (g *game) AddPlayer(playerUUID player.Player) error {
	return nil
}
