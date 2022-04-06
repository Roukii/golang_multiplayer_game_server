package player_action

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

type MoveAction struct {
	Position   entity.Position
	PlayerUUID string
	Created    time.Time
}

func (action MoveAction) Perform(game *game.GameService) {
	player, ok := game.ConnectedPlayer[action.PlayerUUID]
	if !ok {
		return
	}
	player.CurrentPosition = action.Position
	game.SendChange(MovePlayerChange{
		PlayerUUID:   player.UUID,
		Position:     action.Position,
	})
}

type MovePlayerChange struct {
	game.PlayerChange
	PlayerUUID string
	Position   entity.Position
}
