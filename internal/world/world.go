package world

import (
	"fmt"
	"log"
	"net"

	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func Run() {
	l := logger.New("logger.level")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetString("port")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// pb.RegisterGreeterServer(s, &entity.Character{})
	s.RegisterService(&Greeter_ServiceDesc, srv)
	l.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
