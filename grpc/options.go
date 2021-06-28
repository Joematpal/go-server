package grpc

type Option interface {
	applyOption(*Server) error
}
