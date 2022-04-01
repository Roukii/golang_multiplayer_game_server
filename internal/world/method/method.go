package method

import (
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
)

type (
	Server struct {
		pb.UnimplementedWorldServer
	}
)
