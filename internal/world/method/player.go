package method

import (
	"context"
	"io"
	"log"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
)

func (c *Server) GetPlayers(ctx context.Context, request bool) (*pb.GetPlayersReply, error) {
	return &pb.GetPlayersReply{
		Player: []*pb.Player{},
	}, nil
}

func (c *Server) Connect(ctx context.Context, request *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	request.GetPlayerUuid()
	return &pb.ConnectResponse{
		Player:        &pb.Player{},
		World:         &pb.World{},
		Chunks:        []*pb.Chunk{},
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (c *Server) Stream(ctx context.Context, requestStream pb.PlayerService_StreamServer) error {
	log.Println("start new server")
	var max int32

	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := requestStream.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}
		action := req.GetAction()
		switch action.(type) {
		case *pb.PlayerStreamRequest_Move:
			resp := moveplayer(req)
			if err := requestStream.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("send new max=%d", max)
		}
		// update max and send it to stream

	}
}

func moveplayer(resp *pb.PlayerStreamRequest) pb.PlayerStreamResponse {
	move := resp.GetMove()
	return pb.PlayerStreamResponse{
		Action: &pb.PlayerStreamResponse_Move{
			Move: &pb.Move{
				Position: &pb.Position{
					Position: &pb.Vector3{
						X: move.Position.Position.X,
						Y: move.Position.Position.Y,
						Z: move.Position.Position.Z,
					},
					Angle: &pb.Vector3{
						X: move.Position.Angle.X,
						Y: move.Position.Angle.Y,
						Z: move.Position.Angle.Z,
					},
				},
			},
		},
	}
}
