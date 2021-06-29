// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package streamer

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// StreamerClient is the client API for Streamer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamerClient interface {
	StreamPoint(ctx context.Context, opts ...grpc.CallOption) (Streamer_StreamPointClient, error)
}

type streamerClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamerClient(cc grpc.ClientConnInterface) StreamerClient {
	return &streamerClient{cc}
}

func (c *streamerClient) StreamPoint(ctx context.Context, opts ...grpc.CallOption) (Streamer_StreamPointClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Streamer_serviceDesc.Streams[0], "/streamer.Streamer/StreamPoint", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamerStreamPointClient{stream}
	return x, nil
}

type Streamer_StreamPointClient interface {
	Send(*Point) error
	CloseAndRecv() (*Status, error)
	grpc.ClientStream
}

type streamerStreamPointClient struct {
	grpc.ClientStream
}

func (x *streamerStreamPointClient) Send(m *Point) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamerStreamPointClient) CloseAndRecv() (*Status, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Status)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamerServer is the server API for Streamer service.
// All implementations must embed UnimplementedStreamerServer
// for forward compatibility
type StreamerServer interface {
	StreamPoint(Streamer_StreamPointServer) error
	mustEmbedUnimplementedStreamerServer()
}

// UnimplementedStreamerServer must be embedded to have forward compatible implementations.
type UnimplementedStreamerServer struct {
}

func (UnimplementedStreamerServer) StreamPoint(Streamer_StreamPointServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamPoint not implemented")
}
func (UnimplementedStreamerServer) mustEmbedUnimplementedStreamerServer() {}

// UnsafeStreamerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamerServer will
// result in compilation errors.
type UnsafeStreamerServer interface {
	mustEmbedUnimplementedStreamerServer()
}

func RegisterStreamerServer(s grpc.ServiceRegistrar, srv StreamerServer) {
	s.RegisterService(&_Streamer_serviceDesc, srv)
}

func _Streamer_StreamPoint_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamerServer).StreamPoint(&streamerStreamPointServer{stream})
}

type Streamer_StreamPointServer interface {
	SendAndClose(*Status) error
	Recv() (*Point, error)
	grpc.ServerStream
}

type streamerStreamPointServer struct {
	grpc.ServerStream
}

func (x *streamerStreamPointServer) SendAndClose(m *Status) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamerStreamPointServer) Recv() (*Point, error) {
	m := new(Point)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Streamer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "streamer.Streamer",
	HandlerType: (*StreamerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamPoint",
			Handler:       _Streamer_StreamPoint_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "streamer.proto",
}
