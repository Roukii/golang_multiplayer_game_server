package dynamic_entity_action

import (
	"log"
	"time"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/dynamic_entity"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

type DisconnectAction struct {
	Message    string
	PlayerUUID string
	Created    time.Time
}

type DisconnectDynamicEntityChange struct {
	dynamic_entity.DynamicEntityChange
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
	log.Println("disconnect : ", action)
	player, ok := game.PlayerService.ConnectedPlayer[action.PlayerUUID]
	if !ok {
		return
	}
	game.PlayerService.DisconnectPlayer(action.PlayerUUID)
	game.SendDynamicEntityChange(DisconnectDynamicEntityChange{
		PlayerUUID: player.UUID,
		Message:    action.Message,
	})
}
