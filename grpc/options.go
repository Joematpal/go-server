package grpc

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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

func WithServerOptions(serverOpts ...grpc.ServerOption) Option {
	return optionApplyFunc(func(s *Server) error {
		s.serverOptions = append(s.serverOptions, serverOpts...)
		return nil
	})
}

func WithGatewayServiceHandlers(serverHandlers ...GatewayServiceHandler) Option {
	return optionApplyFunc(func(s *Server) error {
		s.gatewayServiceHandlers = append(s.gatewayServiceHandlers, serverHandlers...)
		return nil
	})
}

func WithGatewayServerMuxOptions(opts ...runtime.ServeMuxOption) Option {
	return optionApplyFunc(func(s *Server) error {
		s.gatewayServerMuxOptions = opts
		return nil
	})
}

func WithGatewayAddr(host, port string) Option {
	return optionApplyFunc(func(s *Server) error {
		s.gwHost = host
		s.gwPort = port
		return nil
	})
}

// WithGatewayDialOptions
// If sending in any dial options please include full dial options, such as client credentials
// See `WithGatewayDialCredentials` for a helper in setting client credentials
func WithGatewayDialOptions(opts ...grpc.DialOption) Option {
	return optionApplyFunc(func(s *Server) error {
		s.gatewayDialOptions = append(s.gatewayDialOptions, opts...)
		return nil
	})
}

func WithGatewayDialCredentials(pubCert, privCert string) Option {
	return optionApplyFunc(func(s *Server) error {
		var creds credentials.TransportCredentials
		var err error
		certs, err := ParseCertificates(pubCert, privCert)
		if err != nil {
			return fmt.Errorf("parse certs: %v", err)
		}
		creds = Credentials(certs)
		s.gatewayDialOptions = append(s.gatewayDialOptions, grpc.WithTransportCredentials(creds))
		return nil
	})
}

func WithVersionPath(versionPath string) Option {
	return optionApplyFunc(func(s *Server) error {
		if versionPath != "" {
			s.versionPath = versionPath
		}
		return nil
	})
}
