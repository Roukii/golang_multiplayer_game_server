package method

import (
	"context"
	"io"
	"log"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/client"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/Roukii/pock_multiplayer/pkg/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ChunkMethod struct {
	pb.UnimplementedChunkServiceServer
	game    *game.GameService
	clients *client.ClientService
}

func (c *ChunkMethod) GetWorlds(ctx context.Context, request *emptypb.Empty) (*pb.GetWorldsResponse, error) {
	worlds := c.game.UniverseService.GetWorlds()
	var pbWorlds []*pb.World
	for _, world := range worlds {
		pbWorlds = append(pbWorlds, helper.WorldTypeToProto(world))
	}
	return &pb.GetWorldsResponse{
		Worlds: pbWorlds,
	}, nil
}

func (c *ChunkMethod) EnterWorld(ctx context.Context, request *pb.EnterWorldRequest) (*pb.EnterWorldResponse, error) {
	world, err := c.game.UniverseService.GetWorld(request.WorldUUID)
	if err != nil {
		log.Println("couldn't find world : ", err)
		return nil, status.Errorf(codes.InvalidArgument, "couldn't find world")
	}
	chunks, err := c.game.UniverseService.LoadChunksFromSpawnPoint(world.SpawnPoints[0])
	if err != nil {
		log.Println("couldn't load chunks : ", err)
		return nil, status.Errorf(codes.InvalidArgument, "couldn't load chunks")
	}
	return &pb.EnterWorldResponse{
		World:         helper.WorldTypeToProto(world),
		Chunks:        helper.ChunksTypeToProto(chunks),
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (c *ChunkMethod) LoadChunk(ctx context.Context, request *pb.LoadChunkRequest) (*pb.LoadChunkResponse, error) {
	return &pb.LoadChunkResponse{
		Chunks:        []*pb.Chunk{},
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (c *ChunkMethod) ChunkStream(stream pb.ChunkService_StreamServer) error {
	var lastRequest *pb.ChunkStreamRequest
	for {
		err := stream.RecvMsg(lastRequest)
		if err == io.EOF {
			return stream.Send(&pb.ChunkStreamResponse{
				Action: &pb.ChunkStreamResponse_AddStaticEntity{
					AddStaticEntity: &pb.AddStaticEntity{},
				},
			})
		}
		if err != nil {
			return err
		}
	}
}
