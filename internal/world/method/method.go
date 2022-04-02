package method

import (
	"sync"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type client struct {
	streamChunkServer  pb.ChunkService_StreamServer
	streamPlayerServer pb.PlayerService_StreamServer
	lastMessage        time.Time
	done               chan error
	player             player.Player
}

func New(s grpc.ServiceRegistrar, gameService *game.GameService) {
	pb.RegisterChunkServiceServer(s, &ChunkMethod{
		game:    gameService,
		clients: map[uuid.UUID]*client{},
		mu:      sync.RWMutex{},
	})
	pb.RegisterPlayerServiceServer(s, &PlayerMethod{
		clients: map[uuid.UUID]*client{},
		game:    gameService,
		mu:      sync.RWMutex{},
	})
}
