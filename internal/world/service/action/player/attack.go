package player_action

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

type AttackAction struct {
	DynamicEntityUUID []string
	StaticEntityUUID  []string
	Angle             entity.Vector3f
	PlayerUUID        string
	Created           time.Time
}

type AttackPlayerChange struct {
	game.PlayerChange
	PlayerUUID        string
	DynamicEntityUUID []string
	StaticEntityUUID  []string
	Angle             entity.Vector3f
}

func SendAttackAction(req *pb.PlayerStreamRequest, game *game.GameService, playerUUID string) {
	attack := req.GetAttack()
	game.PlayerActionChannel <- AttackAction{
		DynamicEntityUUID: attack.DynamicEntityUUID,
		StaticEntityUUID:  attack.StaticEntityUUID,
		Angle:             entity.Vector3f{
			X: attack.Angle.X,
			Y: attack.Angle.Y,
			Z: attack.Angle.Z,
		},
		PlayerUUID:        playerUUID,
		Created:           time.Now(),
	}
}

// TODO check if player can hit target
func (action AttackAction) Perform(game *game.GameService) {
	player, ok := game.PlayerService.ConnectedPlayer[action.PlayerUUID]
	if !ok {
		return
	}
	player.CurrentPosition.Rotation = action.Angle
	game.SendPlayerChange(AttackPlayerChange{
		PlayerUUID:        player.UUID,
		DynamicEntityUUID: []string{},
		StaticEntityUUID:  []string{},
		Angle:             action.Angle,
	})
}
