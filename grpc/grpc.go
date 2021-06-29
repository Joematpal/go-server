package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	host             string
	port             string
	tls              bool
	pubCert          string
	privCert         string
	listener         net.Listener
	serverOptions    []grpc.ServerOption
	grpcServer       *grpc.Server
	httpServer       *http.Server
	logger           logger
	registerServices []RegisterService
}

type logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type RegisterService = func(*grpc.Server)

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
		if creds, err := ParseCredentials(s.pubCert, s.privCert); err != nil {
			return nil, err
		} else {
			s.serverOptions = append(s.serverOptions, grpc.Creds(creds))
		}
	}

	s.grpcServer = grpc.NewServer(s.serverOptions...)

	for _, service := range s.registerServices {
		service(s.grpcServer)
	}

	// TODO: Register the http server for the GRPC Gateway.

	// TODO:  and Register the grpc client on the mux runtime; pbSomething.RegisterSomerthingHander

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

	if srv.grpcServer != nil {
		eg.Go(func() error {
			return srv.grpcServer.Serve(srv.listener)
		})
	}

	// before we run the gateway server we need to check if we even need it.
	if srv.httpServer != nil {
		eg.Go(func() error {
			return srv.httpServer.ListenAndServe()
		})
	}

	eg.Go(func() error {
		<-ctx.Done()
		fmt.Println("is it done?")
		if srv.grpcServer != nil {
			srv.grpcServer.GracefulStop()
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
