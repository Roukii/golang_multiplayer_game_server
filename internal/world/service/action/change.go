package action

import (
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	dynamic_entity_action "github.com/Roukii/pock_multiplayer/internal/world/service/action/dynamic_entity"
	"github.com/Roukii/pock_multiplayer/internal/world/service/dynamic_entity"
	"github.com/Roukii/pock_multiplayer/pkg/helper"
)

func GetDynamicEntityChangeToProto(change dynamic_entity.DynamicEntityChange) *pb.PlayerStreamResponse {
	var resp *pb.PlayerStreamResponse
	switch change.(type) {
	case dynamic_entity_action.MoveDynamicEntityChange:
		change := change.(dynamic_entity_action.MoveDynamicEntityChange)
		resp = helper.PlayerMoveChangeToProto(change)
	case dynamic_entity_action.AttackDynamicEntityChange:
		change := change.(dynamic_entity_action.AttackDynamicEntityChange)
		resp = helper.PlayerAttackChangeToProto(change)
	case dynamic_entity_action.ConnectDynamicEntityChange:
		change := change.(dynamic_entity_action.ConnectDynamicEntityChange)
		resp = helper.PlayerConnectChangeToProto(change)
	case dynamic_entity_action.DisconnectDynamicEntityChange:
		change := change.(dynamic_entity_action.DisconnectDynamicEntityChange)
		resp = helper.PlayerDisconnectChangeToProto(change)
	case dynamic_entity_action.HitDynamicEntityChange:
		change := change.(dynamic_entity_action.HitDynamicEntityChange)
		resp = helper.PlayerHitChangeToProto(change)
	case dynamic_entity_action.InteractDynamicEntityChange:
		change := change.(dynamic_entity_action.InteractDynamicEntityChange)
		resp = helper.PlayerInteractChangeToProto(change)
	case dynamic_entity_action.UseSkillDynamicEntityChange:
		change := change.(dynamic_entity_action.UseSkillDynamicEntityChange)
		resp = helper.PlayerUseSkillChangeToProto(change)
	}
	return resp
}
