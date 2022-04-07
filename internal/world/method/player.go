package method

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/client"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	player_action "github.com/Roukii/pock_multiplayer/internal/world/service/game/action/player"
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

	players, err := pm.game.PlayerService.PlayerDao.GetAllPlayersFromUserUUID(userInfo.UUID)
	var playerResponse []*pb.Player
	for _, player := range players {
		playerResponse = append(playerResponse, &pb.Player{
			Name:  player.Name,
			Level: int32(player.Stats.Level),
			Position: &pb.Position{
				Position: &pb.Vector3{
					X: player.SpawnPoint.Coordinate.Position.X,
					Y: player.SpawnPoint.Coordinate.Position.Y,
					Z: player.SpawnPoint.Coordinate.Position.Z,
				},
				Angle: &pb.Vector3{},
			},
		})
	}
	return &pb.GetPlayersReply{
		Player: playerResponse,
	}, nil
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
	world, chunks, err := pm.loadWorldAndChunksFromSpawnPoint(p.SpawnPoint)
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
	world, chunks, err := pm.loadWorldAndChunksFromSpawnPoint(p.SpawnPoint)
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

// TODO lock write with mutex
func (pm *PlayerMethod) loadWorldAndChunksFromSpawnPoint(spawnPoint player.SpawnPoint) (world *universe.World, chunks []*universe.Chunk, err error) {
	world, err = pm.game.UniverseService.GetWorld(spawnPoint.WorldUUID)
	if err != nil {
		fmt.Println("failed to load world", err)
		return nil, nil, status.Errorf(codes.InvalidArgument, "failed to load world")
	}
	chunks, err = pm.game.UniverseService.GetChunksFromSpawnSpoint(spawnPoint, 1)
	if err != nil {
		fmt.Println("failed to load chunks", err)
		return nil, nil, status.Errorf(codes.InvalidArgument, "failed to load chunks")
	}
	return world, chunks, nil
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
			action := req.GetAction()
			switch action.(type) {
			case *pb.PlayerStreamRequest_Move:
				pm.handleMoveAction(req, currentClient.GetPlayerUUID())
			}
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
			switch change.(type) {
			case player_action.MovePlayerChange:
				change := change.(player_action.MovePlayerChange)
				pm.handleMoveChange(change)
			}
		}
	}()
}

func (pm *PlayerMethod) handleMoveAction(req *pb.PlayerStreamRequest, playerUUID string) {
	move := req.GetMove()
	pm.game.PlayerActionChannel <- player_action.MoveAction{
		Position: entity.Position{
			Position: entity.Vector3f{
				X: move.Position.Position.X,
				Y: move.Position.Position.Y,
				Z: move.Position.Position.Z,
			},
			Rotation: entity.Vector3f{
				X: 0,
				Y: 0,
				Z: 0,
			},
		},
		PlayerUUID: playerUUID,
		Created:    time.Now(),
	}
}

func (pm *PlayerMethod) handleMoveChange(move player_action.MovePlayerChange) {
	resp := pb.PlayerStreamResponse{
		Action: &pb.PlayerStreamResponse_Move{
			Move: &pb.Move{
				Position: &pb.Position{
					Position: &pb.Vector3{
						X: move.Position.Position.X,
						Y: move.Position.Position.Y,
						Z: move.Position.Position.Z,
					},
					Angle: &pb.Vector3{
						X: move.Position.Rotation.X,
						Y: move.Position.Rotation.Y,
						Z: move.Position.Rotation.Z,
					},
				},
			},
		},
	}
	pm.clients.BroadcastPlayer(&resp)
}

func handleMovePlayerChange(resp *pb.PlayerStreamRequest) pb.PlayerStreamResponse {
	move := resp.GetMove()
	return pb.PlayerStreamResponse{
		Action: &pb.PlayerStreamResponse_Move{Move: &pb.Move{Position: &pb.Position{Position: &pb.Vector3{X: move.Position.Position.X, Y: move.Position.Position.Y, Z: move.Position.Position.Z}, Angle: &pb.Vector3{X: move.Position.Angle.X, Y: move.Position.Angle.Y, Z: move.Position.Angle.Z}}}},
	}
}
