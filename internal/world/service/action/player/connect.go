package player_action

import (
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

type ConnectPlayerChange struct {
	game.PlayerChange
	Player     *player.Player
}
