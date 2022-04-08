package helper

import (
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	dynamic_entity_action "github.com/Roukii/pock_multiplayer/internal/world/service/action/dynamic_entity"
)

// TODO modify these query to allow creatures instead of only players

func PlayerMoveChangeToProto(move dynamic_entity_action.MoveDynamicEntityChange) *pb.PlayerStreamResponse {
	resp := pb.PlayerStreamResponse{
		Uuid: move.PlayerUUID,
		Info: &pb.PlayerStreamResponse_DynamicEntity{
			DynamicEntity: pb.DynamicEntityType_PLAYER,
		},
		Action: &pb.PlayerStreamResponse_Move{Move: &pb.Move{Position: &pb.Position{Position: &pb.Vector3{X: move.Position.Position.X, Y: move.Position.Position.Y, Z: move.Position.Position.Z}, Angle: &pb.Vector3{X: move.Position.Rotation.X, Y: move.Position.Rotation.Y, Z: move.Position.Rotation.Z}}, Jump: move.Jump}},
	}
	return &resp
}

func PlayerAttackChangeToProto(attack dynamic_entity_action.AttackDynamicEntityChange) *pb.PlayerStreamResponse {
	resp := pb.PlayerStreamResponse{
		Uuid: attack.PlayerUUID,
		Info: &pb.PlayerStreamResponse_DynamicEntity{
			DynamicEntity: pb.DynamicEntityType_PLAYER,
		},
		Action: &pb.PlayerStreamResponse_Attack{Attack: &pb.Attack{DynamicEntityUUID: attack.DynamicEntityUUID, StaticEntityUUID: attack.StaticEntityUUID, Angle: vector3fToProto(attack.Angle)}},
	}
	return &resp
}

func PlayerHitChangeToProto(hit dynamic_entity_action.HitDynamicEntityChange) *pb.PlayerStreamResponse {
	resp := pb.PlayerStreamResponse{
		Uuid: hit.PlayerUUID,
		Info: &pb.PlayerStreamResponse_DynamicEntity{
			DynamicEntity: pb.DynamicEntityType_PLAYER,
		},
		Action: &pb.PlayerStreamResponse_Hit{Hit: &pb.Hit{Damage: hit.Damage, HpLeft: hit.HpLeft, Position: vector3fToProto(hit.Position), DynamicEntityUUID: hit.DynamicEntityUUID, StaticEntityUUID: hit.StaticEntityUUID, SkillId: hit.SkillId}},
	}
	return &resp
}

func PlayerConnectChangeToProto(connect dynamic_entity_action.ConnectDynamicEntityChange) *pb.PlayerStreamResponse {
	resp := pb.PlayerStreamResponse{
		Uuid: connect.Player.UUID,
		Info: &pb.PlayerStreamResponse_DynamicEntity{DynamicEntity: pb.DynamicEntityType_PLAYER},
		Action: &pb.PlayerStreamResponse_Connect{
			Connect: &pb.PlayerConnect{
				Player: &pb.Player{
					Name:  connect.Player.Name,
					Uuid:  connect.Player.UUID,
					Level: int32(connect.Player.Stats.Level),
					Position: &pb.Position{
						Position: vector3fToProto(connect.Player.Position.Position),
						Angle:    vector3fToProto(connect.Player.Position.Rotation),
					},
				},
			},
		},
	}
	return &resp
}

func PlayerDisconnectChangeToProto(disconnect dynamic_entity_action.DisconnectDynamicEntityChange) *pb.PlayerStreamResponse {
	resp := pb.PlayerStreamResponse{
		Uuid: disconnect.PlayerUUID,
		Info: &pb.PlayerStreamResponse_DynamicEntity{DynamicEntity: pb.DynamicEntityType_PLAYER},
		Action: &pb.PlayerStreamResponse_Disconnect{
			Disconnect: &pb.PlayerDisconnect{
				Message: disconnect.Message,
			},
		},
	}
	return &resp
}

func PlayerInteractChangeToProto(interact dynamic_entity_action.InteractDynamicEntityChange) *pb.PlayerStreamResponse {
	var pbInteract pb.Interact
	if interact.DynamicEntityType != nil {
		pbInteract = pb.Interact{
			Uuid: interact.EntityUUID,
			Info: &pb.Interact_DynamicEntity{
				DynamicEntity: pbInteract.GetDynamicEntity(),
			},
		}
	} else if interact.StaticEntityType != nil {
		pbInteract = pb.Interact{
			Uuid: interact.EntityUUID,
			Info: &pb.Interact_StaticEntity{
				StaticEntity: pbInteract.GetStaticEntity(),
			},
		}
	}
	resp := pb.PlayerStreamResponse{
		Uuid: interact.PlayerUUID,
		Info: &pb.PlayerStreamResponse_DynamicEntity{DynamicEntity: pb.DynamicEntityType_PLAYER},
		Action: &pb.PlayerStreamResponse_Interact{
			Interact: &pbInteract,
		},
	}
	return &resp
}

// TODO send skill info
func PlayerUseSkillChangeToProto(skill dynamic_entity_action.UseSkillDynamicEntityChange) *pb.PlayerStreamResponse {
	resp := pb.PlayerStreamResponse{
		Uuid: skill.PlayerUUID,
		Info: &pb.PlayerStreamResponse_DynamicEntity{DynamicEntity: pb.DynamicEntityType_PLAYER},
		Action: &pb.PlayerStreamResponse_Skill{
			Skill: &pb.UseSkill{
				Position: vector3fToProto(skill.Position),
				Angle:    vector3fToProto(skill.Angle),
				Skill: &pb.Skill{
					SkillUuid: skill.SkillId,
				},
				Id: skill.SkillId,
			},
		},
	}
	return &resp
}
