package player_action

import (
	"log"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

// TODO implement with skill
type UseSkillAction struct {
	Position entity.Vector3f
	Angle    entity.Vector3f
	SkillId  string
	// Skill      pb.Skill
	PlayerUUID string
	Created    time.Time
}

type UseSkillPlayerChange struct {
	game.PlayerChange
	PlayerUUID string
	Position   entity.Vector3f
	Angle      entity.Vector3f
	SkillId    string
	// Skill      pb.Skill
}

func SendUseSkillAction(req *pb.PlayerStreamRequest, game *game.GameService, playerUUID string) {
	skill := req.GetSkill()
	game.PlayerActionChannel <- UseSkillAction{
		Position:   entity.Vector3f{X: skill.Position.X, Y: skill.Position.Y, Z: skill.Position.Z},
		Angle:      entity.Vector3f{X: skill.Angle.X, Y: skill.Angle.Y, Z: skill.Angle.Z},
		SkillId:    skill.Skill.SkillUuid,
		PlayerUUID: playerUUID,
		Created:    time.Now(),
	}
}

// TODO check if player can hit target
func (action UseSkillAction) Perform(game *game.GameService) {
	log.Println("%s - use skill : %v", action.PlayerUUID, action)
	player, ok := game.PlayerService.ConnectedPlayer[action.PlayerUUID]
	if !ok {
		return
	}
	player.CurrentPosition.Rotation = action.Angle
	game.SendPlayerChange(UseSkillPlayerChange{
		PlayerUUID: player.UUID,
		Position:   action.Position,
		Angle:      action.Angle,
		SkillId:    action.SkillId,
	})
}
