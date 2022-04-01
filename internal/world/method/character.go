package method

import (
	"context"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
)

func (c *Server) Add(ctx context.Context, request *pb.AddRequest) (*pb.AddReply, error) {
	return &pb.AddReply{
		Message: "gg",
	}, nil
}
