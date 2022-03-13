package world

import (
	"log"
	"net"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWorldServer
}

func Run() {
	l := logger.New("logger.level")

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWorldServer(s, &server{})
	l.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		l.Fatal("failed to serve: %v", err)
	}
}
