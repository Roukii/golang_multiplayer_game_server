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
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	player_action "github.com/Roukii/pock_multiplayer/internal/world/service/game/action/player"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"github.com/Roukii/pock_multiplayer/pkg/helper"

)

const (
	clientTimeout = 15
	maxClients    = 8
)

type PlayerMethod struct {
	pb.UnimplementedPlayerServiceServer
	clients map[string]*client
	game    *game.GameService
	mu      sync.RWMutex
}

func (pm *PlayerMethod) GetPlayers(ctx context.Context, request *emptypb.Empty) (*pb.GetPlayersReply, error) {
	userInfo, err := getUserInfoFromRequest(ctx)
	if err != nil {
		return nil, err
	}

	players, err := pm.game.PlayerDao.GetAllPlayersFromUserUUID(userInfo.UUID)
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
	pm.game.PlayerMu.Lock()
	err = pm.game.CreatePlayer(userInfo.UUID, &p)
	pm.game.PlayerMu.Unlock()
	if err != nil {
		fmt.Println("failed to create player", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to create player")
	}
	world, chunks, err := pm.loadWorldAndChunksFromSpawnPoint(p.SpawnPoint)
	if err != nil {
		return nil, err
	}

	pm.mu.Lock()
	pm.clients[userInfo.UUID] = &client{
		lastMessage: time.Now(),
		done:        make(chan error),
		playerUUID:  p.UUID,
		userUUID:  userInfo.UUID,
	}
	pm.mu.Unlock()

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
	pm.game.PlayerMu.Lock()
	player_uuid := request.GetPlayerUuid()
	p, err := pm.game.ConnectPlayer(userInfo.UUID, player_uuid)
	pm.game.PlayerMu.Unlock()
	if err != nil {
		fmt.Println("failed to connect player", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to create player")
	}
	world, chunks, err := pm.loadWorldAndChunksFromSpawnPoint(p.SpawnPoint)
	if err != nil {
		return nil, err
	}

	pm.mu.Lock()
	pm.clients[p.UUID] = &client{
		lastMessage: time.Now(),
		done:        make(chan error),
		playerUUID:  p.UUID,
		userUUID:  userInfo.UUID,
	}
	pm.mu.Unlock()

	return &pb.ConnectResponse{
		Player:        helper.PlayerTypeToProto(p),
		World:         helper.WorldTypeToProto(world),
		Chunks:        helper.ChunksTypeToProto(chunks),
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

// TODO lock write with mutex
func (pm *PlayerMethod) loadWorldAndChunksFromSpawnPoint(spawnPoint player.SpawnPoint) (world *universe.World, chunks []*universe.Chunk, err error) {
	world, err = pm.game.GetWorld(spawnPoint.WorldUUID)
	if err != nil {
		fmt.Println("failed to load world", err)
		return nil, nil, status.Errorf(codes.InvalidArgument, "failed to load world")
	}
	chunks, err = pm.game.GetChunksFromSpawnSpoint(spawnPoint, 1)
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
	currentClient, ok := pm.clients[userInfo.UUID]
	if !ok {
		fmt.Println("client no found")
		return status.Errorf(codes.InvalidArgument, "user not recognized")
	}
	if currentClient.streamPlayerServer != nil {
		return errors.New("stream already active")
	}
	currentClient.streamPlayerServer = requestStream
	fmt.Println("Start new stream for : ", currentClient.playerUUID)
	go func() {
		for {
			// receive data from stream
			req, err := requestStream.Recv()
			if err != nil {
				log.Printf("receive error %v", err)
				currentClient.done <- errors.New("failed to receive request")
				return
			}
			action := req.GetAction()
			switch action.(type) {
			case *pb.PlayerStreamRequest_Move:
				pm.handleMoveAction(req, currentClient)
			}
		}
	}()

	// Wait for stream to be done.
	var doneError error
	select {
	case <-ctx.Done():
		doneError = ctx.Err()
	case doneError = <-currentClient.done:
	}
	log.Printf(`stream done with error "%v"`, doneError)

	log.Printf("%s - removing client", currentClient.playerUUID)
	pm.removeClient(currentClient.userUUID)
	pm.removePlayer(currentClient.playerUUID)

	return doneError
}

func (pm *PlayerMethod) removeClient(userUUID string) {
	pm.mu.Lock()
	delete(pm.clients, userUUID)
	pm.mu.Unlock()
}

func (pm *PlayerMethod) removePlayer(playerUUID string) {
	pm.game.PlayerMu.Lock()
	pm.game.DisconnectPlayer(playerUUID)
	pm.game.PlayerMu.Unlock()

	// TODO broadcast disconnect
	// resp := proto.Response{
	// 	Action: &proto.Response_RemoveEntity{
	// 		RemoveEntity: &proto.RemoveEntity{
	// 			Id: playerID.String(),
	// 		},
	// 	},
	// }
	// s.broadcast(&resp)
}

func (pm *PlayerMethod) watchTimeout() {
	timeoutTicker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			for _, client := range pm.clients {
				if time.Now().Sub(client.lastMessage).Minutes() > clientTimeout {
					client.done <- errors.New("you have been timed out")
					return
				}
			}
			<-timeoutTicker.C
		}
	}()
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

func (pm *PlayerMethod) handleMoveAction(req *pb.PlayerStreamRequest, currentClient *client) {
	move := req.GetMove()
	pm.game.PlayerActionChannel <- player_action.MoveAction{
		Position:   entity.Position{
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
		PlayerUUID: currentClient.playerUUID,
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
	pm.broadcast(&resp)
}

func handleMovePlayerChange(resp *pb.PlayerStreamRequest) pb.PlayerStreamResponse {
	move := resp.GetMove()
	return pb.PlayerStreamResponse{
		Action: &pb.PlayerStreamResponse_Move{Move: &pb.Move{Position: &pb.Position{Position: &pb.Vector3{X: move.Position.Position.X, Y: move.Position.Position.Y, Z: move.Position.Position.Z}, Angle: &pb.Vector3{X: move.Position.Angle.X, Y: move.Position.Angle.Y, Z: move.Position.Angle.Z}}}},
	}
}

func (pm *PlayerMethod) broadcast(resp *pb.PlayerStreamResponse) {
	pm.mu.Lock()
	for id, currentClient := range pm.clients {
		if currentClient.streamPlayerServer == nil {
			continue
		}
		if err := currentClient.streamPlayerServer.Send(resp); err != nil {
			log.Printf("%s - broadcast error %v", id, err)
			currentClient.done <- errors.New("failed to broadcast message")
			continue
		}
		log.Printf("%s - broadcasted %+v", resp, id)
	}
	pm.mu.Unlock()
}
