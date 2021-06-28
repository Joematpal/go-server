package grpc

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"

	"github.com/digital-dream-labs/go-server/flags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func New(opts ...Option) error {
	// set up grpc connection
	// TODO: Iterate over opts and apply them to option struct
	// TODO: needs to be option instead of c.String

	grpcHost := c.String(flags.GRPCHost)
	grpcPort := c.String(flags.GRPCPort)
	isTLS := c.Bool(flags.GRPCTLS)
	grpcPubCert := c.String(flags.GRPCPubCert)
	grpcPrivCert := c.String(flags.GRPCPrivCert)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", grpcHost, grpcPort))
	if err != nil {
		return err
	}

	// TODO: Add to option struct
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpclog.UnaryServerInterceptor()),
		grpc.StreamInterceptor(grpclog.StreamServerInterceptor()),
	}

	if isTLS {
		// log.Infoln("tls is on")
		creds, err := credentials.NewServerTLSFromFile(grpcPubCert, grpcPrivCert)
		if err != nil {
			pubCert, err := base64.StdEncoding.DecodeString(grpcPubCert)
			if err != nil {
				return err
			}
			privCert, err := base64.StdEncoding.DecodeString(grpcPrivCert)
			if err != nil {
				return err
			}
			creds, err := tls.X509KeyPair(pubCert, privCert)
			if err != nil {
				return err
			}
			opts = append(opts, grpc.Creds(credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{creds}})))
		}
		opts = append(opts, grpc.Creds(creds))
	}
	srv := grpc.NewServer(opts...)

	// TODO: Register here

	fmt.Println("listening:", grpcHost+":"+grpcPort)
	return srv.Serve(lis)
}
