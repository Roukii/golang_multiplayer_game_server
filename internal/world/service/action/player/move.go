package player_action

import (
	"log"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

type MoveAction struct {
	Position   entity.Position
	Jump       bool
	PlayerUUID string
	Created    time.Time
}

type MovePlayerChange struct {
	game.PlayerChange
	PlayerUUID string
	Position   entity.Position
	Jump       bool
}

func SendMoveAction(req *pb.PlayerStreamRequest, game *game.GameService, playerUUID string) {
	move := req.GetMove()
	game.PlayerActionChannel <- MoveAction{
		Position:   entity.Position{Position: entity.Vector3f{X: move.Position.Position.X, Y: move.Position.Position.Y, Z: move.Position.Position.Z}, Rotation: entity.Vector3f{X: move.Position.Angle.X, Y: move.Position.Angle.X, Z: move.Position.Angle.X}},
		Jump:       move.Jump,
		PlayerUUID: playerUUID,
		Created:    time.Now(),
	}
}

func (action MoveAction) Perform(game *game.GameService) {
	log.Println("move : ", action)
	player, ok := game.PlayerService.ConnectedPlayer[action.PlayerUUID]
	if !ok {
		return
	}
	player.CurrentPosition = action.Position
	game.SendPlayerChange(MovePlayerChange{
		PlayerUUID: player.UUID,
		Position:   action.Position,
		Jump:       action.Jump,
	})
}
