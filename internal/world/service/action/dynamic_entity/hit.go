package dynamic_entity_action

import (
	"log"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/dynamic_entity"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

type HitAction struct {
	DynamicEntityUUID []string
	StaticEntityUUID  []string
	Damage            int64
	HpLeft            int64
	Position          entity.Vector3f
	SkillId           string
	PlayerUUID        string
	Created           time.Time
}

type HitDynamicEntityChange struct {
	dynamic_entity.DynamicEntityChange
	PlayerUUID        string
	DynamicEntityUUID []string
	StaticEntityUUID  []string
	Damage            int64
	HpLeft            int64
	Position          entity.Vector3f
	SkillId           string
}

func SendHitAction(req *pb.PlayerStreamRequest, game *game.GameService, playerUUID string) {
	hit := req.GetHit()
	game.PlayerActionChannel <- HitAction{
		DynamicEntityUUID: hit.DynamicEntityUUID,
		StaticEntityUUID:  hit.StaticEntityUUID,
		Damage:            hit.Damage,
		HpLeft:            hit.HpLeft,
		Position:          entity.Vector3f{X: hit.Position.X, Y: hit.Position.Y, Z: hit.Position.Z},
		SkillId:           hit.SkillId,
		PlayerUUID:        playerUUID,
		Created:           time.Now(),
	}
}

// TODO check if player can hit target also have array of damage instead of just a single value
func (action HitAction) Perform(game *game.GameService) {
	log.Println("hit : ", action)
	player, ok := game.PlayerService.ConnectedPlayer[action.PlayerUUID]
	if !ok {
		return
	}
	game.SendDynamicEntityChange(HitDynamicEntityChange{
		PlayerUUID:        player.UUID,
		DynamicEntityUUID: []string{},
		StaticEntityUUID:  []string{},
		Damage:            action.Damage,
		HpLeft:            action.Damage,
		Position:          action.Position,
		SkillId:           action.SkillId,
	})
}
