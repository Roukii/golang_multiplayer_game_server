package method

import (
	"context"
	"strings"
	"sync"
	"time"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/Roukii/pock_multiplayer/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type client struct {
	streamChunkServer  pb.ChunkService_StreamServer
	streamPlayerServer pb.PlayerService_StreamServer
	lastMessage        time.Time
	done               chan error
	userUUID           string
	playerUUID         string
}

func New(s grpc.ServiceRegistrar, gameService *game.GameService) {
	chunkServer := ChunkMethod{
		game:    gameService,
		clients: map[string]*client{},
		mu:      sync.RWMutex{},
	}
	playerServer := PlayerMethod{
		clients: map[string]*client{},
		game:    gameService,
		mu:      sync.RWMutex{},
	}

	pb.RegisterChunkServiceServer(s, &chunkServer)
	pb.RegisterPlayerServiceServer(s, &playerServer)
	playerServer.watchChanges()
}

func getUserInfoFromRequest(ctx context.Context) (*jwt.CustomerInfo, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "failed to get metadata")
	}
	token := md["token"]
	if len(token) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing 'token' header")
	}
	if strings.Trim(token[0], " ") == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty 'token' header")
	}
	userInfo, err := jwt.VerifyToken(token[0])
	if err != nil {
		return &jwt.CustomerInfo{
			UUID:   token[0],
			Name:   "",
			Device: "",
		}, nil
		// return nil, status.Errorf(codes.InvalidArgument, "invalid token")
	}
	return userInfo, nil
}
