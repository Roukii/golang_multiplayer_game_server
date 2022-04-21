package helper

import (
	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
)

// TODO update chunk width and population with real value
func WorldTypeToProto(world *universe.World) *pb.World {
	return &pb.World{
		Uuid:        world.UUID,
		Name:        world.Name,
		Level:       int32(world.Level),
		Length:      int32(world.Length),
		Width:       int32(world.Width),
		ScaleXY:     int32(world.ScaleXY),
		ScaleHeight: int32(world.ScaleHeight),
		Seed:        world.Seed,
		ChunkWidth:  25,
		Population:  int32(world.MaxPlayer),
	}
}

func PlayerTypeToProto(player *player.Player) *pb.Player {
	return &pb.Player{
		DynamicEntity: DynamicEntityToProto(&player.IDynamicEntity),
	}
}

func ChunksTypeToProto(chunks []*universe.Chunk) []*pb.Chunk {
	var pbChunks []*pb.Chunk
	for _, chunk := range chunks {
		var tiles []*pb.Tile
		for _, tile := range chunk.Tiles {
			tiles = append(tiles, &pb.Tile{
				Type:      pb.TileType(tile.TileType),
				Elevation: float32(tile.Elevation),
			})
		}
		pbChunks = append(pbChunks, &pb.Chunk{
			Uuid:         chunk.UUID,
			Position:     &pb.Vector2Int{X: int32(chunk.PositionX), Y: int32(chunk.PositionY)},
			StaticEntity: []*pb.StaticEntity{},
			Tiles:        tiles,
		})
	}
	return pbChunks
}

func DynamicEntitiesToProto(dynamicEntities map[string]entity.DynamicEntity) []*pb.DynamicEntity {
	var pbDynamicEntities []*pb.DynamicEntity
	for _, de := range dynamicEntities {
		pbDynamicEntities = append(pbDynamicEntities, DynamicEntityToProto(de))
	}
	return pbDynamicEntities
}

func DynamicEntityToProto(de entity.DynamicEntity) *pb.DynamicEntity {
	return &pb.DynamicEntity{
		Uuid: de.GetUUID(),
		Name: de.GetName(),
		Position: &pb.Position{
			Position: Vector3fToProto(de.GetPosition().Position),
			Angle:    Vector3fToProto(de.GetPosition().Rotation),
		},
		Type: pb.DynamicEntityType(de.GetType()),
		Stats: &pb.Stats{
			Level: int64(de.GetStats().Level),
			MaxHP: int64(de.GetStats().Maxhp),
			HP:    int64(de.GetStats().Hp),
			MaxMP: int64(de.GetStats().Maxmp),
			MP:    int64(de.GetStats().Mp),
		},
	}
}

func Vector3fToProto(pos entity.Vector3f) *pb.Vector3 {
	return &pb.Vector3{
		X: pos.X,
		Y: pos.Y,
		Z: pos.Z,
	}
}
