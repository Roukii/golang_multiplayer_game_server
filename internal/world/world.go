package world

import (
	"log"
	"net"

	"github.com/Roukii/pock_multiplayer/internal/world/method"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/gocql/gocql"
	"google.golang.org/grpc"
)

func Run() {
	l := logger.New("logger.level")

	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWorldServer(s, &method.Server{})
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Keyspace = "game"
	session, err := cluster.CreateSession()
	if err != nil {
		l.Fatal("failed to create session: %v", err)
	}
	service.New(session, l)
	l.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		l.Fatal("failed to serve: %v", err)
	}
}
