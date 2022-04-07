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
		Name:        world.Name,
		Level:       int32(world.Level),
		Length:      int32(world.Length),
		Width:       int32(world.Width),
		ScaleXY:     1,
		ScaleHeight: 8,
		Seed:        world.Seed,
		ChunkWidth:  25,
		Population:  int32(world.MaxPlayer),
	}
}

func PlayerTypeToProto(player *player.Player) *pb.Player {
	return &pb.Player{
		Name:  player.Name,
		Uuid:  player.UUID,
		Level: int32(player.Stats.Level),
		Position: &pb.Position{
			Position: vector3fToProto(player.CurrentPosition.Position),
			Angle: vector3fToProto(player.CurrentPosition.Rotation),
		},
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

func vector3fToProto(pos entity.Vector3f) *pb.Vector3 {
	return &pb.Vector3{
		X: pos.X,
		Y: pos.Y,
		Z: pos.Z,
	}
}