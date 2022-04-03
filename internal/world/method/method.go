package method

import (
	"context"
	"sync"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/Roukii/pock_multiplayer/pkg/jwt"
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

func getUserInfoFromRequest(ctx context.Context) (*jwt.CustomerInfo, error) {
	return &jwt.CustomerInfo{
		UUID:   "7e7067e0-b2d2-11ec-885d-367dda4cfa8c",
		Name:   "",
		Device: "",
	}, nil
	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	return nil, status.Errorf(codes.DataLoss, "failed to get metadata")
	// }
	// token := md["token"]
	// if len(token) == 0 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "missing 'token' header")
	// }
	// if strings.Trim(token[0], " ") == "" {
	// 	return nil, status.Errorf(codes.InvalidArgument, "empty 'token' header")
	// }
	// userInfo, err := jwt.VerifyToken(token[0])
	// if err != nil {
	// 	return nil, status.Errorf(codes.InvalidArgument, "invalid token")
	// }
	// return userInfo, nil
}
