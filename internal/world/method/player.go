package method

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/action"
	"github.com/Roukii/pock_multiplayer/internal/world/service/client"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/Roukii/pock_multiplayer/pkg/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type PlayerMethod struct {
	pb.UnimplementedPlayerServiceServer
	clients *client.ClientService
	game    *game.GameService
	mu      sync.RWMutex
}

func (pm *PlayerMethod) GetPlayers(ctx context.Context, request *emptypb.Empty) (*pb.GetPlayersReply, error) {
	userInfo, err := getUserInfoFromRequest(ctx)
	if err != nil {
		return nil, err
	}
	playersReply, err := pm.game.PlayerService.GetPlayersFromUserUUID(userInfo.UUID)
	return playersReply, nil
}

func (pm *PlayerMethod) CreatePlayer(ctx context.Context, request *pb.CreatePlayerRequest) (*pb.CreatePlayerResponse, error) {
	userInfo, err := getUserInfoFromRequest(ctx)
	if err != nil {
		return nil, err
	}
	p := player.Player{
		Name: request.GetName(),
		Stats: entity.Stats{
			Level: 1,
			Maxhp: 10,
			Hp:    10,
			Maxmp: 10,
			Mp:    10,
		},
	}
	// TODO need to transform this into spawn point
	var worldUUID string
	for _, world := range pm.game.UniverseService.Universe.Worlds {
		worldUUID = world.UUID
		break
	}
	err = pm.game.PlayerService.CreatePlayer(userInfo.UUID, &p, worldUUID)
	if err != nil {
		fmt.Println("failed to create player", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to create player")
	}
	world, chunks, err := pm.game.UniverseService.LoadWorldAndChunksFromSpawnPoint(p.SpawnPoint)
	if err != nil {
		return nil, err
	}

	pm.clients.AddClient(userInfo.UUID, &p)

	return &pb.CreatePlayerResponse{
		Player:        helper.PlayerTypeToProto(&p),
		World:         helper.WorldTypeToProto(world),
		Chunks:        helper.ChunksTypeToProto(chunks),
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (pm *PlayerMethod) Connect(ctx context.Context, request *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	userInfo, err := getUserInfoFromRequest(ctx)
	if err != nil {
		return nil, err
	}
	player_uuid := request.GetPlayerUuid()
	p, err := pm.game.PlayerService.ConnectPlayer(userInfo.UUID, player_uuid)
	if err != nil {
		fmt.Println("failed to connect player", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to create player")
	}
	world, chunks, err := pm.game.UniverseService.LoadWorldAndChunksFromSpawnPoint(p.SpawnPoint)
	if err != nil {
		return nil, err
	}

	pm.clients.AddClient(userInfo.UUID, p)

	return &pb.ConnectResponse{
		Player:        helper.PlayerTypeToProto(p),
		World:         helper.WorldTypeToProto(world),
		Chunks:        helper.ChunksTypeToProto(chunks),
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (pm *PlayerMethod) Stream(requestStream pb.PlayerService_StreamServer) error {
	ctx := requestStream.Context()
	userInfo, err := getUserInfoFromRequest(ctx)
	if err != nil {
		return err
	}
	currentClient, ok := pm.clients.GetClient(userInfo.UUID)
	if !ok {
		fmt.Println("client no found")
		return status.Errorf(codes.InvalidArgument, "user not recognized")
	}

	err = currentClient.AddPlayerStream(requestStream)
	if err != nil {
		return errors.New("stream already active")
	}

	go func() {
		for {
			// receive data from stream
			req, err := requestStream.Recv()
			if err != nil {
				log.Printf("receive error %v", err)
				currentClient.Done <- errors.New("failed to receive request")
				return
			}
			action.SendPlayerAction(req, pm.game, currentClient.GetPlayerUUID())
		}
	}()

	// Wait for stream to be done.
	var doneError error
	select {
	case <-ctx.Done():
		doneError = ctx.Err()
	case doneError = <-currentClient.Done:
	}
	log.Printf(`stream done with error "%v"`, doneError)

	pm.clients.DisconnectClient(currentClient, doneError.Error())

	return doneError
}

func (pm *PlayerMethod) watchChanges() {
	go func() {
		for {
			change := <-pm.game.PlayerChangeChannel
			pm.clients.BroadcastPlayer(action.GetPlayerChangeToProto(change))
		}
	}()
}
