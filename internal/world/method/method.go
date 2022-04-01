package method

import (
	"sync"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/google/uuid"
)

type (
	Server struct {
		pb.UnimplementedChunkServiceServer
		pb.UnimplementedPlayerServiceServer
		clients map[uuid.UUID]*client
		game    *game.Game
		mu      sync.RWMutex
	}
)

type client struct {
	streamChunkServer  pb.ChunkService_StreamServer
	streamPlayerServer pb.PlayerService_StreamServer
	lastMessage        time.Time
	done               chan error
	player             player.Player
}
