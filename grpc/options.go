package grpc

type Option interface {
	applyOption(*Options) error
}

type Options struct {
}
