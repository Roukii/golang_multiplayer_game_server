package action

import (
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	player_action "github.com/Roukii/pock_multiplayer/internal/world/service/action/player"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/Roukii/pock_multiplayer/pkg/helper"
)

func GetPlayerChangeToProto(change game.PlayerChange) *pb.PlayerStreamResponse {
	var resp *pb.PlayerStreamResponse
	switch change.(type) {
	case player_action.MovePlayerChange:
		change := change.(player_action.MovePlayerChange)
		resp = helper.PlayerMoveChangeToProto(change)
	case player_action.AttackPlayerChange:
		change := change.(player_action.AttackPlayerChange)
		resp = helper.PlayerAttackChangeToProto(change)
	case player_action.ConnectPlayerChange:
		change := change.(player_action.ConnectPlayerChange)
		resp = helper.PlayerConnectChangeToProto(change)
	case player_action.DisconnectPlayerChange:
		change := change.(player_action.DisconnectPlayerChange)
		resp = helper.PlayerDisconnectChangeToProto(change)
	case player_action.HitPlayerChange:
		change := change.(player_action.HitPlayerChange)
		resp = helper.PlayerHitChangeToProto(change)
	case player_action.InteractPlayerChange:
		change := change.(player_action.InteractPlayerChange)
		resp = helper.PlayerInteractChangeToProto(change)
	case player_action.UseSkillPlayerChange:
		change := change.(player_action.UseSkillPlayerChange)
		resp = helper.PlayerUseSkillChangeToProto(change)
	}
	return resp
}
