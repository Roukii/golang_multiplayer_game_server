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
 // TODO add token check everywhere
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

// TODO check if world can be entered and move player dynamic entity between world service
func (c *ChunkMethod) EnterWorld(ctx context.Context, request *pb.EnterWorldRequest) (*pb.EnterWorldResponse, error) {
	userInfo, err := getUserInfoFromRequest(ctx)
	if err != nil {
		return nil, err
	}
	client, ok := c.clients.GetClient(userInfo.UUID)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "player not connected")
	}
	player, ok := c.game.PlayerService.ConnectedPlayer[client.GetPlayerUUID()]
	if !ok {
		log.Println("player not found")
		return nil, status.Errorf(codes.NotFound, "player not found")
	}

	world, err := c.game.UniverseService.GetWorldService(request.WorldUUID)
	if err != nil {
		log.Println("couldn't find world : ", err)
		return nil, status.Errorf(codes.InvalidArgument, "couldn't find world")
	}
	// TODO better handle of spawn point please
	chunks, err := world.LoadChunksFromSpawnPoint(world.World.SpawnPoints[0])
	if err != nil {
		log.Println("couldn't load chunks : ", err)
		return nil, status.Errorf(codes.InvalidArgument, "couldn't load chunks")
	}
	player.CurrentWorldUUID = world.World.UUID
	return &pb.EnterWorldResponse{
		World:         helper.WorldTypeToProto(world.World),
		Chunks:        helper.ChunksTypeToProto(chunks),
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (c *ChunkMethod) LoadChunk(ctx context.Context, request *pb.LoadChunkRequest) (*pb.LoadChunkResponse, error) {
	userInfo, err := getUserInfoFromRequest(ctx)
	if err != nil {
		return nil, err
	}
	client, ok := c.clients.GetClient(userInfo.UUID)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "player not connected")
	}
	player, ok := c.game.PlayerService.ConnectedPlayer[client.GetPlayerUUID()]
	if !ok {
		log.Println("player not found")
		return nil, status.Errorf(codes.NotFound, "player not found")
	}
	worldService, err := c.game.UniverseService.GetWorldService(player.CurrentWorldUUID)
	if err != nil {
		return nil, err
	}
	chunks, err := worldService.LoadSpecificChunks(request.ChunkToLoad)
	if err != nil {
		return nil, err
	}
	return &pb.LoadChunkResponse{
		Chunks:        helper.ChunksTypeToProto(chunks),
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (c *ChunkMethod) ChunkStream(stream pb.ChunkService_StreamServer) error {
	var lastRequest *pb.ChunkStreamRequest
	for {
		err := stream.RecvMsg(lastRequest)
		if err == io.EOF {
			return stream.Send(&pb.ChunkStreamResponse{
				AddStaticEntity:    []*pb.AddStaticEntity{},
				UpdateStaticEntity: []*pb.UpdateStaticEntity{},
				RemoveStaticEntity: []*pb.RemoveStaticEntity{},
			})
		}
	}
}
