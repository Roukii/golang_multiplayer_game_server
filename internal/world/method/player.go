package method

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
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
	client, ok := pm.clients.GetClient(userInfo.UUID)
	if ok && client.GetPlayerUUID() != "" {
		return nil, status.Errorf(codes.AlreadyExists, "already connect")
	}

	p := player.Player{
		IDynamicEntity: entity.IDynamicEntity{Name: request.GetName(), Stats: entity.Stats{Level: 1, Maxhp: 10, Hp: 10, Maxmp: 10, Mp: 10}},
	}
	// TODO choose a world in another way
	var world *universe.World
	for _, w := range pm.game.UniverseService.WorldServices {
		world = w.World
		log.Println(w.World)
		log.Println(w.World.SpawnPoints)
		break
	}
	log.Println("create player")
	err = pm.game.PlayerService.CreatePlayer(userInfo.UUID, &p, world)
	log.Println("finished create player")
	if err != nil {
		log.Println("failed to create player", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to create player")
	}
	worldService, err := pm.game.UniverseService.GetWorldService(p.CurrentWorldUUID)
	if err != nil {
		log.Println("can't load world : ", err)
		return nil, err
	}
	chunks, err := worldService.LoadChunksFromSpawnPoint(p.SpawnPoint)
	if err != nil {
		log.Println("can´t load chunk : ", err)
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
	client, ok := pm.clients.GetClient(userInfo.UUID)
	if ok && client.GetPlayerUUID() != "" {
		return nil, status.Errorf(codes.AlreadyExists, "already connect")
	}
	player_uuid := request.GetPlayerUuid()
	p, err := pm.game.PlayerService.ConnectPlayer(userInfo.UUID, player_uuid)
	if err != nil {
		log.Println("failed to connect player", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to create player")
	}
	worldService, err := pm.game.UniverseService.GetWorldService(p.SpawnPoint.WorldUUID)
	chunks, err := worldService.LoadChunksFromSpawnPoint(p.SpawnPoint)
	if err != nil {
		return nil, err
	}

	pm.clients.AddClient(userInfo.UUID, p)
	return &pb.ConnectResponse{
		Player:        helper.PlayerTypeToProto(p),
		World:         helper.WorldTypeToProto(worldService.World),
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
		log.Println("client not found")
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
				currentClient.PlayerDone <- errors.New("failed to receive request")
				return
			}
			currentClient.Update()
			action.SendPlayerAction(req, pm.game, currentClient.GetPlayerUUID())
		}
	}()

	// Wait for stream to be done.
	var doneError error
	select {
	case <-ctx.Done():
		doneError = ctx.Err()
	case doneError = <-currentClient.PlayerDone:
	}
	log.Printf(`stream done with error "%v"`, doneError)
	currentClient.RemovePlayerStream()

	return doneError
}

func (pm *PlayerMethod) watchChanges() {
	go func() {
		for {
			change := <-pm.game.DynamicEntityChangeChannel
			pm.clients.BroadcastPlayer(action.GetDynamicEntityChangeToProto(change))
		}
	}()
}
