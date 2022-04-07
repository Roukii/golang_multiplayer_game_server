package action

import (
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	player_action "github.com/Roukii/pock_multiplayer/internal/world/service/action/player"

)

func SendPlayerAction(req *pb.PlayerStreamRequest, game *game.GameService, playerUUID string) {
	action := req.GetAction()
	switch action.(type) {
	case *pb.PlayerStreamRequest_Move:
		player_action.SendMoveAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Attack:
		player_action.SendAttackAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Disconnect:
		player_action.SendDisconnectAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Hit:
		player_action.SendHitAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Interact:
		player_action.SendInteractAction(req, game, playerUUID)
	case *pb.PlayerStreamRequest_Skill:
		player_action.SendUseSkillAction(req, game, playerUUID)
	}
}
