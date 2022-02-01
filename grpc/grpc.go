package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Handler interface {
	Handle(path string, handler http.Handler)
	http.Handler
}
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

	versionPath string

	// TLS Option for ignoring the tls in
	insecureSkipVerify bool
	dialCerts          []tls.Certificate
	clientAuthType     tls.ClientAuthType
	clientCAs          *x509.CertPool

	// Gateway
	gwConn                  *grpc.ClientConn
	gwHost                  string
	gwPort                  string
	gatewayServiceHandlers  []GatewayServiceHandler
	gatewayServerMuxOptions []runtime.ServeMuxOption
	gatewayDialOptions      []grpc.DialOption

	// Experimental
	handler Handler

	// Swagger
	swaggerFile string
}

type logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type RegisterService = func(*grpc.Server)

type GatewayServiceHandler = func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) (err error)

// set up grpc connection
// Takes in a anonymous function that registers the service
func New(opts ...Option) (*Server, error) {
	s := &Server{
		serverOptions:           []grpc.ServerOption{},
		registerServices:        []RegisterService{},
		gatewayServerMuxOptions: []runtime.ServeMuxOption{},
		gatewayServiceHandlers:  []GatewayServiceHandler{},
		gatewayDialOptions:      []grpc.DialOption{},
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

	//
	if s.IsTLS() && (s.dialCerts == nil || s.dialCerts != nil && len(s.dialCerts) == 0) {
		certs, err := ParseCertificates(s.pubCert, s.privCert)
		if err != nil {
			return nil, fmt.Errorf("parse certs: %w", err)
		}
		s.dialCerts = []tls.Certificate{certs}
	}

	return s, nil
}

// Takes the receiver and copies them to the new server
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

	if len(s.gatewayServerMuxOptions) > 0 {
		server.gatewayServerMuxOptions = append(server.gatewayServerMuxOptions, s.gatewayServerMuxOptions...)
	}

	if len(s.gatewayServiceHandlers) > 0 {
		server.gatewayServiceHandlers = append(server.gatewayServiceHandlers, s.gatewayServiceHandlers...)
	}

	if len(s.gatewayDialOptions) > 0 {
		server.gatewayDialOptions = append(server.gatewayDialOptions, s.gatewayDialOptions...)
	}

	// This shouldn't ever need to be copied, but it probably better to be safe than sorry
	if s.gwConn != nil {
		server.gwConn = s.gwConn
	}

	if s.gwHost != "" {
		server.gwHost = s.gwHost
	}

	if s.gwPort != "" {
		server.gwPort = s.gwPort
	}

	if s.swaggerFile != "" {
		server.swaggerFile = s.swaggerFile
	}

	if s.versionPath != "" {
		server.versionPath = s.versionPath
	}

	return nil
}

func (s *Server) IsTLS() bool {
	return s.tls
}

func (s *Server) StartWithContext(ctx context.Context) error {
	// this will be difficult because of it will need to handle two different go routines for spinning off the grpc server and the http gateway server

	if s.IsTLS() {
		s.serverOptions = append(s.serverOptions, grpc.Creds(credentials.NewTLS(
			&tls.Config{
				Certificates:       s.dialCerts,
				ClientCAs:          s.getCertPool(),
				ClientAuth:         s.clientAuthType,
				InsecureSkipVerify: s.insecureSkipVerify,
			},
		)))
	}

	if len(s.serverOptions) != 0 {
		s.grpcServer = grpc.NewServer(s.serverOptions...)
	}

	for _, service := range s.registerServices {
		service(s.grpcServer)
	}

	// Creates GRPC gateway if needed
	if len(s.gatewayServiceHandlers) > 0 {
		s.Debugf("running gRPC gateway")
		if err := s.newGRPCGateway(ctx); err != nil {
			return fmt.Errorf("new grpc gateway: %v", err)
		}
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run grpc server only when gateway is not running - because httpServer has a mux to grpc
	// and dont want two listeners on same port
	if s.grpcServer != nil && s.httpServer == nil {
		s.Debugf("running gRPC at %s", s.port)
		eg.Go(func() error {
			return s.grpcServer.Serve(s.listener)
		})
	}

	// before we run the gateway server we need to check if we even need it.
	if s.httpServer != nil {
		eg.Go(func() error {
			s.Infof("http listening at %s", s.httpServer.Addr)
			if s.IsTLS() {
				return s.httpServer.Serve(tls.NewListener(s.listener, s.httpServer.TLSConfig))
			}
			return s.httpServer.Serve(s.listener)
		})
	}

	eg.Go(func() error {
		<-ctx.Done()
		s.Debugf("start shutdown")
		if s.grpcServer != nil {
			s.grpcServer.GracefulStop()
		}
		if s.httpServer != nil {
			if err := s.httpServer.Shutdown(ctx); err != nil {
				s.Debugf("http shutdown: %v", err)
			}
		}
		if s.gwConn != nil {
			if err := s.gwConn.Close(); err != nil {
				s.Debugf("gateway client conn: %v", err)
			}
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
