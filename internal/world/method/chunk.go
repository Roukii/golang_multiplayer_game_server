package method

import (
	"context"
	"io"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
)

func (c *Server) EnterChunk(ctx context.Context, request *pb.EnterChunkRequest) (*pb.EnterChunkResponse, error) {
	return &pb.EnterChunkResponse{
		Chunks:        []*pb.Chunk{},
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (c *Server) LoadChunk(ctx context.Context, request *pb.LoadChunkRequest) (*pb.LoadChunkResponse, error) {
	return &pb.LoadChunkResponse{
		Chunks:        []*pb.Chunk{},
		DynamicEntity: []*pb.DynamicEntity{},
	}, nil
}

func (c *Server) ChunkStream(stream pb.ChunkService_StreamServer) error {
	var lastRequest *pb.ChunkStreamRequest
	for {
		err := stream.RecvMsg(lastRequest)
		if err == io.EOF {
			return stream.Send(&pb.ChunkStreamResponse{
				Action: &pb.ChunkStreamResponse_AddStaticEntity{
					AddStaticEntity: &pb.AddStaticEntity{},
				},
			})
		}
		if err != nil {
			return err
		}
	}
}
