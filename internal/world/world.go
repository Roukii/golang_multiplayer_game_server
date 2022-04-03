package world

import (
	"log"
	"net"

	"github.com/Roukii/pock_multiplayer/internal/world/method"
	"github.com/Roukii/pock_multiplayer/internal/world/service"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"google.golang.org/grpc"
)

func Run() {
	l := logger.New("logger.level")

	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	cluster := gocql.NewCluster("127.0.0.1:9042")
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		l.Fatal("create session:", err)
	}
	err = session.ExecStmt(`CREATE KEYSPACE IF NOT EXISTS game WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`)
	if err != nil {
		l.Fatal("create keyspace:", err)
	}
	service := service.New(&session, l)
	service.GameService.StartGame()
	method.New(s, service.GameService)
	l.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		l.Fatal("failed to serve: %v", err)
	}
}
