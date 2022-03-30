package world

import (
	"context"
	"log"
	"net"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
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
	pb.RegisterWorldServer(s, &entity.Server{})
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Keyspace = "game"
	session, err := cluster.CreateSession()
	if err != nil {
		l.Fatal("failed to create session: %v", err)
	}

	err = session.Query("create index on game.tweet(timeline)").WithContext(context.Background()).Exec()
	if err != nil {
		l.Fatal("failed to INSERT: %v", err)
	}

	err = session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").WithContext(context.Background()).Exec()
	if err != nil {
		l.Fatal("failed to INSERT: %v", err)
	}

	var id string
	var text string
	err = session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`,
		"me").WithContext(context.Background()).Consistency(gocql.One).Scan(&id, &text)
	if err != nil {
		l.Fatal("failed to SELECT: %v", err)
	}

	l.Debug(id, text)
	l.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		l.Fatal("failed to serve: %v", err)
	}
}
