package game

import "github.com/Roukii/pock_multiplayer/internal/world/service/dynamic_entity"

type PlayerAction interface {
	Perform(game *GameService)
}

func (game *GameService) SendDynamicEntityChange(change dynamic_entity.DynamicEntityChange) {
	select {
	case game.DynamicEntityChangeChannel <- change:
	default:
	}
}
