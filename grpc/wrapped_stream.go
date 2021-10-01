package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type WrappedStream struct {
	ctx context.Context
	grpc.ServerStream
}

func (ws *WrappedStream) Context() context.Context {
	return ws.ctx
}

func NewWrappedStream(ctx context.Context, s grpc.ServerStream) grpc.ServerStream {
	return &WrappedStream{ctx, s}
}
