package action

import (
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	dynamic_entity_action "github.com/Roukii/pock_multiplayer/internal/world/service/action/dynamic_entity"
	"github.com/Roukii/pock_multiplayer/internal/world/service/dynamic_entity"
	"github.com/Roukii/pock_multiplayer/pkg/helper/proto_action"
)

func GetDynamicEntityChangeToProto(change dynamic_entity.DynamicEntityChange) *pb.PlayerStreamResponse {
	var resp *pb.PlayerStreamResponse
	switch change.(type) {
	case dynamic_entity_action.MoveDynamicEntityChange:
		change := change.(dynamic_entity_action.MoveDynamicEntityChange)
		resp = proto_action.PlayerMoveChangeToProto(change)
	case dynamic_entity_action.AttackDynamicEntityChange:
		change := change.(dynamic_entity_action.AttackDynamicEntityChange)
		resp = proto_action.PlayerAttackChangeToProto(change)
	case dynamic_entity_action.ConnectDynamicEntityChange:
		change := change.(dynamic_entity_action.ConnectDynamicEntityChange)
		resp = proto_action.PlayerConnectChangeToProto(change)
	case dynamic_entity_action.DisconnectDynamicEntityChange:
		change := change.(dynamic_entity_action.DisconnectDynamicEntityChange)
		resp = proto_action.PlayerDisconnectChangeToProto(change)
	case dynamic_entity_action.HitDynamicEntityChange:
		change := change.(dynamic_entity_action.HitDynamicEntityChange)
		resp = proto_action.PlayerHitChangeToProto(change)
	case dynamic_entity_action.InteractDynamicEntityChange:
		change := change.(dynamic_entity_action.InteractDynamicEntityChange)
		resp = proto_action.PlayerInteractChangeToProto(change)
	case dynamic_entity_action.UseSkillDynamicEntityChange:
		change := change.(dynamic_entity_action.UseSkillDynamicEntityChange)
		resp = proto_action.PlayerUseSkillChangeToProto(change)
	}
	return resp
}
