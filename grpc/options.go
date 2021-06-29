package grpc

import "google.golang.org/grpc"

type Option interface {
	applyOption(*Server) error
}

type optionApplyFunc func(*Server) error

func (f optionApplyFunc) applyOption(opts *Server) error {
	return f(opts)
}

func WithRegisterService(rs RegisterService) Option {
	return optionApplyFunc(func(s *Server) error {
		s.registerServices = append(s.registerServices, rs)
		return nil
	})
}

func WithHost(host string) Option {
	return optionApplyFunc(func(s *Server) error {
		s.host = host
		return nil
	})
}

func WithPort(port string) Option {
	return optionApplyFunc(func(s *Server) error {
		s.port = port
		return nil
	})
}

func WithTLS(tls bool) Option {
	return optionApplyFunc(func(s *Server) error {
		s.tls = tls
		return nil
	})
}

func WithPubCert(pubCert string) Option {
	return optionApplyFunc(func(s *Server) error {
		s.pubCert = pubCert
		return nil
	})
}

func WithPrivCert(privCert string) Option {
	return optionApplyFunc(func(s *Server) error {
		s.privCert = privCert
		return nil
	})
}

func WithLogger(logr logger) Option {
	return optionApplyFunc(func(s *Server) error {
		s.logger = logr
		return nil
	})
}

func WithServerOptions(serverOpts []grpc.ServerOption) Option {
	return optionApplyFunc(func(s *Server) error {
		s.serverOptions = append(s.serverOptions, serverOpts...)
		return nil
	})
}
