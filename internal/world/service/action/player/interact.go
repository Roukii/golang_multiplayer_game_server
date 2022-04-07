package player_action

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

type InteractAction struct {
	PlayerUUID        string
	Created           time.Time
	EntityUUID        string
	DynamicEntityType *entity.DynamicEntityType
	StaticEntityType  *entity.StaticEntityType
}

type InteractPlayerChange struct {
	game.PlayerChange
	PlayerUUID        string
	EntityUUID        string
	DynamicEntityType *entity.DynamicEntityType
	StaticEntityType  *entity.StaticEntityType
}

func SendInteractAction(req *pb.PlayerStreamRequest, game *game.GameService, playerUUID string) {
	interact := req.GetInteract()
	interactAction := InteractAction{
		PlayerUUID:        playerUUID,
		Created:           time.Now(),
		EntityUUID:        interact.Uuid,
		DynamicEntityType: nil,
		StaticEntityType:  nil,

	}
	switch interact.GetInfo().(type) {
		case *pb.Interact_DynamicEntity:
			dynamicEntityType := entity.DynamicEntityType(interact.GetDynamicEntity())
			interactAction.DynamicEntityType = &dynamicEntityType
		case *pb.Interact_StaticEntity:
			staticEntityType := entity.StaticEntityType(interact.GetStaticEntity())
			interactAction.StaticEntityType = &staticEntityType
	}
	game.PlayerActionChannel <- interactAction
}

// TODO check if player can interact with target
func (action InteractAction) Perform(game *game.GameService) {
	player, ok := game.PlayerService.ConnectedPlayer[action.PlayerUUID]
	if !ok {
		return
	}
	game.SendPlayerChange(InteractPlayerChange{
		PlayerUUID:        player.UUID,
		EntityUUID:        action.EntityUUID,
		DynamicEntityType: action.DynamicEntityType,
		StaticEntityType:  action.StaticEntityType,
	})
}
