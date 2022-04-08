package action

import (
	"errors"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	dynamic_entity_action "github.com/Roukii/pock_multiplayer/internal/world/service/action/dynamic_entity"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

func SendPlayerAction(req *pb.PlayerStreamRequest, game *game.GameService, playerUUID string) (bool, error) {
	action := req.GetAction()
	switch action.(type) {
	case *pb.PlayerStreamRequest_Move:
		dynamic_entity_action.SendMoveAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Attack:
		dynamic_entity_action.SendAttackAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Hit:
		dynamic_entity_action.SendHitAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Interact:
		dynamic_entity_action.SendInteractAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Skill:
		dynamic_entity_action.SendUseSkillAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Disconnect:
		dynamic_entity_action.SendDisconnectAction(req, game, playerUUID)
		return true, errors.New("disconnect")
	}
	return false, nil
}
