package dynamic_entity_action

import (
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/service/dynamic_entity"
)

type ConnectDynamicEntityChange struct {
	dynamic_entity.DynamicEntityChange
	Player *player.Player
}
