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
	host          string
	port          string
	tls           bool
	listener      net.Listener
	serverOptions []grpc.ServerOption
	grpcServer    *grpc.Server
	httpServer    *http.Server
}

// set up grpc connection
func New(opts ...Option) (*Server, error) {
	// TODO: Iterate over o.serverOptions and apply them to option struct
	o := &Server{}
	var err error
	for _, opt := range opts {
		if err := opt.applyOption(o); err != nil {
			return nil, err
		}
	}
	// TODO: needs to be option instead of c.String

	// grpcHost := c.String(flags.GRPCHost)
	// grpcPort := c.String(flags.GRPCPort)
	// isTLS := c.Bool(flags.GRPCTLS)
	// grpcPubCert := c.String(flags.GRPCPubCert)
	// grpcPrivCert := c.String(flags.GRPCPrivCert)

	if o.listener == nil {
		o.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", o.host, o.port))
		if err != nil {
			return nil, err
		}
	}

	// TODO: Add to option struct
	// o.serverOptions := []grpc.ServerOption{
	// 	grpc.UnaryInterceptor(grpclog.UnaryServerInterceptor()),
	// 	grpc.StreamInterceptor(grpclog.StreamServerInterceptor()),
	// }

	if o.IsTLS() {
		// log.Infoln("tls is on")
		// creds, err := credentials.NewServerTLSFromFile(grpcPubCert, grpcPrivCert)
		// if err != nil {
		// 	pubCert, err := base64.StdEncoding.DecodeString(grpcPubCert)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	privCert, err := base64.StdEncoding.DecodeString(grpcPrivCert)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	creds, err := tls.X509KeyPair(pubCert, privCert)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	o.serverOptions = append(o.serverOptions, grpc.Creds(credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{creds}})))
		// }
		// o.serverOptions = append(o.serverOptions, grpc.Creds(creds))
	}

	o.grpcServer = grpc.NewServer(o.serverOptions...)

	// TODO: Register here grpc interface to the implementation

	// TODO: Register the http server for the GRPC Gateway.
	// and Register the grpc client on the mux runtime; pbSomething.RegisterSomerthingHander

	fmt.Println("listening:", o.host+":"+o.port)
	return o, nil
}

// Copies the struct
func (o *Server) applyOption(*Server) error {
	// TODO: copy struct
	return nil
}

func (o *Server) IsTLS() bool {
	return o.tls
}

func (srv *Server) StartWithContext(ctx context.Context) error {
	// this will be difficult because of it iwll need to handle two different go routines for spinning off the grpc server and the http gateway server

	eg := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return srv.grpcServer.Serve(srv.listener)
	})

	// before we run the gateway server we need to check if we even need it.
	eg.Go(func() error {
		return srv.httpServer.ListenAndServe()
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			//do a check if we have a httpServer
			// src.httpServer.Shutdown(context.Background())
			srv.grpcServer.GracefulStop()
			return ctx.Err()
		}
	})

	return eg.Wait()
}

// TODO : think about canceling
