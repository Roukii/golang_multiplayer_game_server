package player_action

import (
	"time"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

type DisconnectAction struct {
	Message    string
	PlayerUUID string
	Created    time.Time
}

type DisconnectPlayerChange struct {
	game.PlayerChange
	PlayerUUID string
	Message    string
}

func SendDisconnectAction(req *pb.PlayerStreamRequest, game *game.GameService, playerUUID string) {
	disconnect := req.GetDisconnect()
	game.PlayerActionChannel <- DisconnectAction{
		Message:    disconnect.Message,
		PlayerUUID: playerUUID,
		Created:    time.Now(),
	}
}

// TODO check if player can hit target
func (action DisconnectAction) Perform(game *game.GameService) {
	player, ok := game.PlayerService.ConnectedPlayer[action.PlayerUUID]
	if !ok {
		return
	}
	game.PlayerService.DisconnectPlayer(action.PlayerUUID)
	game.SendPlayerChange(DisconnectPlayerChange{
		PlayerUUID:   player.UUID,
		Message:      action.Message,
	})
}
