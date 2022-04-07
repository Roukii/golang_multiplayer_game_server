package player_service

import (
	"sync"

	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type PlayerService struct {
	ConnectedPlayer map[string]player.Player
	Logger          *logger.Logger
	Mu              sync.RWMutex
	playerDao       *dao.PlayerDao
}

func NewPlayerService(session *gocqlx.Session) *PlayerService {
	return &PlayerService{
		ConnectedPlayer: make(map[string]player.Player),
		playerDao:       dao.NewPlayerDao(session),
		Logger:          &logger.Logger{},
		Mu:              sync.RWMutex{},
	}
}

func (ps *PlayerService) CreatePlayer(userUuid string, p *player.Player, world *universe.World) (err error) {
	p.UUID = gocql.TimeUUID().String()
	p.SpawnPoint = world.SpawnPoints[0]
	if err != nil {
		return err
	}
	p.CurrentPosition = p.SpawnPoint.Coordinate
	err = ps.playerDao.Insert(userUuid, p)
	if err != nil {
		return err
	}
	ps.Mu.Lock()
	ps.ConnectedPlayer[p.UUID] = *p
	ps.Mu.Unlock()
	return nil
}

func (ps *PlayerService) ConnectPlayer(userUUID string, playerUUID string) (*player.Player, error) {
	p, err := ps.playerDao.GetPlayerFromUUID(userUUID, playerUUID)
	if err != nil {
		return nil, err
	}
	ps.Mu.Lock()
	ps.ConnectedPlayer[p.UUID] = *p
	ps.Mu.Unlock()
	return p, err
}

func (ps *PlayerService) DisconnectPlayer(playerUUID string) (bool, error) {
	if p, ok := ps.ConnectedPlayer[playerUUID]; ok {
		ps.Mu.Lock()
		delete(ps.ConnectedPlayer, p.UUID)
		ps.Mu.Unlock()
		err := ps.playerDao.Update(&p)
		if err != nil {
			return false, err
		}
		return true, err
	}
	return false, nil
}

func (ps *PlayerService) GetPlayersFromUserUUID(userUUID string) (*pb.GetPlayersReply, error) {
	players, err := ps.playerDao.GetAllPlayersFromUserUUID(userUUID)
	if err != nil {
		return nil, err
	}
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
