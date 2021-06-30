package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	host                   string
	port                   string
	tls                    bool
	pubCert                string
	privCert               string
	listener               net.Listener
	serverOptions          []grpc.ServerOption
	grpcServer             *grpc.Server
	httpServer             *http.Server
	logger                 logger
	registerServices       []RegisterService
	gatewayServiceHandlers []GatewayServiceHandler
}

type logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type RegisterService = func(*grpc.Server)
type GatewayServiceHandler = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// set up grpc connection
// Takes in a anonymous function that registers the service
func New(opts ...Option) (*Server, error) {
	s := &Server{
		serverOptions:    []grpc.ServerOption{},
		registerServices: []RegisterService{},
	}
	var err error
	for _, opt := range opts {
		if err := opt.applyOption(s); err != nil {
			return nil, err
		}
	}

	if s.listener == nil {
		s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
		if err != nil {
			return nil, err
		}
	}

	if s.IsTLS() {
		certs, err := ParseCertificates(s.pubCert, s.privCert)
		if err != nil {
			return nil, fmt.Errorf("parse certs: %v", err)
		}
		s.serverOptions = append(s.serverOptions, grpc.Creds(Credentials(certs)))
	}

	s.grpcServer = grpc.NewServer(s.serverOptions...)

	for _, service := range s.registerServices {
		service(s.grpcServer)
	}

	// Creates GRPC gateway if needed
	if len(s.gatewayServiceHandlers) > 0 && s.IsTLS() {
		s.Debugf("running gRPC gateway")
		if err := s.newGRPCGateway(); err != nil {
			return nil, err
		}
	}

	s.Infof("listening: %s:%s", s.host, s.port)
	return s, nil
}

// Takes the receiver and and applies it to the server
func (s *Server) applyOption(server *Server) error {
	if s.host != "" {
		server.host = s.host
	}
	if s.port != "" {
		server.port = s.port
	}
	if s.tls {
		server.tls = s.tls

		if s.pubCert != "" {
			server.pubCert = s.pubCert
		}
		if s.privCert != "" {
			server.privCert = s.privCert
		}
	}
	if s.listener != nil {
		server.listener = s.listener
	}
	if len(s.serverOptions) > 0 {
		server.serverOptions = append(server.serverOptions, s.serverOptions...)
	}

	return nil
}

func (s *Server) IsTLS() bool {
	return s.tls
}

func (srv *Server) StartWithContext(ctx context.Context) error {
	// this will be difficult because of it will need to handle two different go routines for spinning off the grpc server and the http gateway server

	eg, ctx := errgroup.WithContext(ctx)

	// Run grpc server only when gateway is not running - because httpServer has a mux to grpc
	// and dont want two listeners on same port
	if srv.grpcServer != nil && srv.httpServer == nil {
		srv.Debugf("running gRPC")
		eg.Go(func() error {
			return srv.grpcServer.Serve(srv.listener)
		})
	}

	// before we run the gateway server we need to check if we even need it.
	if srv.httpServer != nil {
		srv.Debugf("running http")
		eg.Go(func() error {
			return srv.httpServer.Serve(tls.NewListener(srv.listener, srv.httpServer.TLSConfig))
		})
	}

	eg.Go(func() error {
		<-ctx.Done()
		srv.Debugf("start shutdown")
		if srv.grpcServer != nil {
			srv.grpcServer.GracefulStop()
		}
		if srv.httpServer != nil {
			srv.httpServer.Shutdown(ctx)
		}
		return ctx.Err()
	})

	return eg.Wait()
}

func (s *Server) Infof(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Infof(format, args...)
	}
}

func (s *Server) Debugf(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Debugf(format, args...)
	}
}
