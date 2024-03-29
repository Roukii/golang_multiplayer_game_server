// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.14.0
// source: proto/world.proto

package universe

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ChunkServiceClient is the client API for ChunkService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChunkServiceClient interface {
	GetWorlds(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetWorldsResponse, error)
	EnterWorld(ctx context.Context, in *EnterWorldRequest, opts ...grpc.CallOption) (*EnterWorldResponse, error)
	LoadChunk(ctx context.Context, in *LoadChunkRequest, opts ...grpc.CallOption) (*LoadChunkResponse, error)
	Stream(ctx context.Context, in *ChunkStreamRequest, opts ...grpc.CallOption) (ChunkService_StreamClient, error)
}

type chunkServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChunkServiceClient(cc grpc.ClientConnInterface) ChunkServiceClient {
	return &chunkServiceClient{cc}
}

func (c *chunkServiceClient) GetWorlds(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetWorldsResponse, error) {
	out := new(GetWorldsResponse)
	err := c.cc.Invoke(ctx, "/universe.ChunkService/GetWorlds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chunkServiceClient) EnterWorld(ctx context.Context, in *EnterWorldRequest, opts ...grpc.CallOption) (*EnterWorldResponse, error) {
	out := new(EnterWorldResponse)
	err := c.cc.Invoke(ctx, "/universe.ChunkService/EnterWorld", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chunkServiceClient) LoadChunk(ctx context.Context, in *LoadChunkRequest, opts ...grpc.CallOption) (*LoadChunkResponse, error) {
	out := new(LoadChunkResponse)
	err := c.cc.Invoke(ctx, "/universe.ChunkService/LoadChunk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chunkServiceClient) Stream(ctx context.Context, in *ChunkStreamRequest, opts ...grpc.CallOption) (ChunkService_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChunkService_ServiceDesc.Streams[0], "/universe.ChunkService/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &chunkServiceStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChunkService_StreamClient interface {
	Recv() (*ChunkStreamResponse, error)
	grpc.ClientStream
}

type chunkServiceStreamClient struct {
	grpc.ClientStream
}

func (x *chunkServiceStreamClient) Recv() (*ChunkStreamResponse, error) {
	m := new(ChunkStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChunkServiceServer is the server API for ChunkService service.
// All implementations must embed UnimplementedChunkServiceServer
// for forward compatibility
type ChunkServiceServer interface {
	GetWorlds(context.Context, *emptypb.Empty) (*GetWorldsResponse, error)
	EnterWorld(context.Context, *EnterWorldRequest) (*EnterWorldResponse, error)
	LoadChunk(context.Context, *LoadChunkRequest) (*LoadChunkResponse, error)
	Stream(*ChunkStreamRequest, ChunkService_StreamServer) error
	mustEmbedUnimplementedChunkServiceServer()
}

// UnimplementedChunkServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChunkServiceServer struct {
}

func (UnimplementedChunkServiceServer) GetWorlds(context.Context, *emptypb.Empty) (*GetWorldsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWorlds not implemented")
}
func (UnimplementedChunkServiceServer) EnterWorld(context.Context, *EnterWorldRequest) (*EnterWorldResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnterWorld not implemented")
}
func (UnimplementedChunkServiceServer) LoadChunk(context.Context, *LoadChunkRequest) (*LoadChunkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadChunk not implemented")
}
func (UnimplementedChunkServiceServer) Stream(*ChunkStreamRequest, ChunkService_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedChunkServiceServer) mustEmbedUnimplementedChunkServiceServer() {}

// UnsafeChunkServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChunkServiceServer will
// result in compilation errors.
type UnsafeChunkServiceServer interface {
	mustEmbedUnimplementedChunkServiceServer()
}

func RegisterChunkServiceServer(s grpc.ServiceRegistrar, srv ChunkServiceServer) {
	s.RegisterService(&ChunkService_ServiceDesc, srv)
}

func _ChunkService_GetWorlds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChunkServiceServer).GetWorlds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/universe.ChunkService/GetWorlds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChunkServiceServer).GetWorlds(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChunkService_EnterWorld_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnterWorldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChunkServiceServer).EnterWorld(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/universe.ChunkService/EnterWorld",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChunkServiceServer).EnterWorld(ctx, req.(*EnterWorldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChunkService_LoadChunk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadChunkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChunkServiceServer).LoadChunk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/universe.ChunkService/LoadChunk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChunkServiceServer).LoadChunk(ctx, req.(*LoadChunkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChunkService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChunkStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChunkServiceServer).Stream(m, &chunkServiceStreamServer{stream})
}

type ChunkService_StreamServer interface {
	Send(*ChunkStreamResponse) error
	grpc.ServerStream
}

type chunkServiceStreamServer struct {
	grpc.ServerStream
}

func (x *chunkServiceStreamServer) Send(m *ChunkStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ChunkService_ServiceDesc is the grpc.ServiceDesc for ChunkService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChunkService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "universe.ChunkService",
	HandlerType: (*ChunkServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWorlds",
			Handler:    _ChunkService_GetWorlds_Handler,
		},
		{
			MethodName: "EnterWorld",
			Handler:    _ChunkService_EnterWorld_Handler,
		},
		{
			MethodName: "LoadChunk",
			Handler:    _ChunkService_LoadChunk_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _ChunkService_Stream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/world.proto",
}

// PlayerServiceClient is the client API for PlayerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlayerServiceClient interface {
	CreatePlayer(ctx context.Context, in *CreatePlayerRequest, opts ...grpc.CallOption) (*CreatePlayerResponse, error)
	GetPlayers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetPlayersReply, error)
	Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectResponse, error)
	Stream(ctx context.Context, opts ...grpc.CallOption) (PlayerService_StreamClient, error)
}

type playerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPlayerServiceClient(cc grpc.ClientConnInterface) PlayerServiceClient {
	return &playerServiceClient{cc}
}

func (c *playerServiceClient) CreatePlayer(ctx context.Context, in *CreatePlayerRequest, opts ...grpc.CallOption) (*CreatePlayerResponse, error) {
	out := new(CreatePlayerResponse)
	err := c.cc.Invoke(ctx, "/universe.PlayerService/CreatePlayer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerServiceClient) GetPlayers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetPlayersReply, error) {
	out := new(GetPlayersReply)
	err := c.cc.Invoke(ctx, "/universe.PlayerService/GetPlayers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerServiceClient) Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectResponse, error) {
	out := new(ConnectResponse)
	err := c.cc.Invoke(ctx, "/universe.PlayerService/Connect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerServiceClient) Stream(ctx context.Context, opts ...grpc.CallOption) (PlayerService_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &PlayerService_ServiceDesc.Streams[0], "/universe.PlayerService/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &playerServiceStreamClient{stream}
	return x, nil
}

type PlayerService_StreamClient interface {
	Send(*PlayerStreamRequest) error
	Recv() (*PlayerStreamResponse, error)
	grpc.ClientStream
}

type playerServiceStreamClient struct {
	grpc.ClientStream
}

func (x *playerServiceStreamClient) Send(m *PlayerStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *playerServiceStreamClient) Recv() (*PlayerStreamResponse, error) {
	m := new(PlayerStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PlayerServiceServer is the server API for PlayerService service.
// All implementations must embed UnimplementedPlayerServiceServer
// for forward compatibility
type PlayerServiceServer interface {
	CreatePlayer(context.Context, *CreatePlayerRequest) (*CreatePlayerResponse, error)
	GetPlayers(context.Context, *emptypb.Empty) (*GetPlayersReply, error)
	Connect(context.Context, *ConnectRequest) (*ConnectResponse, error)
	Stream(PlayerService_StreamServer) error
	mustEmbedUnimplementedPlayerServiceServer()
}

// UnimplementedPlayerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPlayerServiceServer struct {
}

func (UnimplementedPlayerServiceServer) CreatePlayer(context.Context, *CreatePlayerRequest) (*CreatePlayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePlayer not implemented")
}
func (UnimplementedPlayerServiceServer) GetPlayers(context.Context, *emptypb.Empty) (*GetPlayersReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlayers not implemented")
}
func (UnimplementedPlayerServiceServer) Connect(context.Context, *ConnectRequest) (*ConnectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedPlayerServiceServer) Stream(PlayerService_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedPlayerServiceServer) mustEmbedUnimplementedPlayerServiceServer() {}

// UnsafePlayerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlayerServiceServer will
// result in compilation errors.
type UnsafePlayerServiceServer interface {
	mustEmbedUnimplementedPlayerServiceServer()
}

func RegisterPlayerServiceServer(s grpc.ServiceRegistrar, srv PlayerServiceServer) {
	s.RegisterService(&PlayerService_ServiceDesc, srv)
}

func _PlayerService_CreatePlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePlayerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServiceServer).CreatePlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/universe.PlayerService/CreatePlayer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServiceServer).CreatePlayer(ctx, req.(*CreatePlayerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerService_GetPlayers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServiceServer).GetPlayers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/universe.PlayerService/GetPlayers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServiceServer).GetPlayers(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerService_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServiceServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/universe.PlayerService/Connect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServiceServer).Connect(ctx, req.(*ConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PlayerServiceServer).Stream(&playerServiceStreamServer{stream})
}

type PlayerService_StreamServer interface {
	Send(*PlayerStreamResponse) error
	Recv() (*PlayerStreamRequest, error)
	grpc.ServerStream
}

type playerServiceStreamServer struct {
	grpc.ServerStream
}

func (x *playerServiceStreamServer) Send(m *PlayerStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *playerServiceStreamServer) Recv() (*PlayerStreamRequest, error) {
	m := new(PlayerStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PlayerService_ServiceDesc is the grpc.ServiceDesc for PlayerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlayerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "universe.PlayerService",
	HandlerType: (*PlayerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePlayer",
			Handler:    _PlayerService_CreatePlayer_Handler,
		},
		{
			MethodName: "GetPlayers",
			Handler:    _PlayerService_GetPlayers_Handler,
		},
		{
			MethodName: "Connect",
			Handler:    _PlayerService_Connect_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _PlayerService_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/world.proto",
}
