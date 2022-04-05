package method

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type PlayerMethod struct {
	pb.UnimplementedPlayerServiceServer
	clients map[uuid.UUID]*client
	game    *game.GameService
	mu      sync.RWMutex
}

func (c *PlayerMethod) GetPlayers(ctx context.Context, request *emptypb.Empty) (*pb.GetPlayersReply, error) {
	userInfo, err := getUserInfoFromRequest(ctx)
	if err != nil {
		return nil, err
	}

	players, err := c.game.PlayerDao.GetAllPlayersFromUserUUID(userInfo.UUID)
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

func (c *PlayerMethod) CreatePlayer(ctx context.Context, request *pb.CreatePlayerRequest) (*pb.CreatePlayerResponse, error) {
	start := time.Now()
	userInfo, err := getUserInfoFromRequest(ctx)
	fmt.Println(userInfo)
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
	err = c.game.CreatePlayer(userInfo.UUID, &p)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to create player")

	}
	elapsed := time.Since(start)
	log.Printf("Load player took %s", elapsed)
	start = time.Now()

	world, err := c.game.GetWorld(p.SpawnPoint.WorldUUID)
	elapsed = time.Since(start)
	log.Printf("Load world took %s", elapsed)
	start = time.Now()
	chunks, err := c.game.GetChunksFromSpawnSpoint(p.SpawnPoint, 1)
	elapsed = time.Since(start)
	log.Printf("Load chunks took %s", elapsed)
	start = time.Now()

	var requestChunk []*pb.Chunk
	for _, chunk := range chunks {
		var tiles []*pb.Tile
		for _, tile := range chunk.Tiles {
			tiles = append(tiles, &pb.Tile{
				Type:      pb.TileType(tile.TileType),
				Elevation: float32(tile.Elevation),
			})
		}
		requestChunk = append(requestChunk, &pb.Chunk{
			Uuid:         chunk.UUID,
			Position:     &pb.Vector2{X: float32(chunk.PositionX), Y: float32(chunk.PositionY)},
			StaticEntity: []*pb.StaticEntity{},
			Tiles:        tiles,
		})
	}
	elapsed = time.Since(start)
	log.Printf("Transfer chunk took %s", elapsed)

	return &pb.CreatePlayerResponse{
		Player:        &pb.Player{Name: p.Name, Level: int32(p.Stats.Level), Position: &pb.Position{Position: &pb.Vector3{X: p.SpawnPoint.Coordinate.Position.X, Y: p.SpawnPoint.Coordinate.Position.Y, Z: p.SpawnPoint.Coordinate.Position.Z}, Angle: &pb.Vector3{}}},
		World:         &pb.World{Name: world.Name, Level: int32(world.Level)},
		Chunks:        requestChunk,
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (c *PlayerMethod) Connect(ctx context.Context, request *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	request.GetPlayerUuid()
	return &pb.ConnectResponse{
		Player:        &pb.Player{},
		World:         &pb.World{},
		Chunks:        []*pb.Chunk{},
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (c *PlayerMethod) Stream(requestStream pb.PlayerService_StreamServer) error {
	log.Println("start new server")
	var max int32
	ctx := requestStream.Context()
	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := requestStream.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}
		action := req.GetAction()
		switch action.(type) {
		case *pb.PlayerStreamRequest_Move:
			resp := moveplayer(req)
			if err := requestStream.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("send new max=%d", max)
		}
		// update max and send it to stream

	}
}

func moveplayer(resp *pb.PlayerStreamRequest) pb.PlayerStreamResponse {
	move := resp.GetMove()
	return pb.PlayerStreamResponse{
		Action: &pb.PlayerStreamResponse_Move{
			Move: &pb.Move{
				Position: &pb.Position{
					Position: &pb.Vector3{
						X: move.Position.Position.X,
						Y: move.Position.Position.Y,
						Z: move.Position.Position.Z,
					},
					Angle: &pb.Vector3{
						X: move.Position.Angle.X,
						Y: move.Position.Angle.Y,
						Z: move.Position.Angle.Z,
					},
				},
			},
		},
	}
}
