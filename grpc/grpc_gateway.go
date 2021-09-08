package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Get Gateway Endpoint check to see if a different Gateway Address is set
func (s *Server) getGatewayEndpoint() string {
	if s.gwHost != "" && s.gwPort != "" {
		return fmt.Sprintf("%s:%s", s.gwHost, s.gwPort)
	}
	return fmt.Sprintf("%s:%s", s.host, s.port)
}

// Register the http server for the GRPC Gateway.
// and Register the grpc client on the mux runtime (handler)
func (s *Server) newGRPCGateway(ctx context.Context) error {

	dialOpts := []grpc.DialOption{}
	// Dial Credentials need to come from the outside if they are different than the local service
	// If no Certificates are pass we assume they are running on the same server
	if len(s.gatewayDialOptions) == 0 {
		var creds credentials.TransportCredentials
		var err error
		certs, err := ParseCertificates(s.pubCert, s.privCert)
		if err != nil {
			return fmt.Errorf("parse certs: %v", err)
		}
		creds = Credentials(certs)
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}

	dialOpts = append(dialOpts, s.gatewayDialOptions...)

	// Register the gateway service handlers; service handlers currently only talk to same grpc.ClientConn
	// (joematpal) I do not see a use case where we need to have maintain a one to many type of connection for the Gateway
	gwmux := runtime.NewServeMux(s.gatewayServerMuxOptions...)
	endpoint := s.getGatewayEndpoint()

	var err error
	s.gwConn, err = grpc.DialContext(ctx, endpoint, dialOpts...)
	if err != nil {
		return fmt.Errorf("dial context: %v", err)
	}

	for _, handler := range s.gatewayServiceHandlers {
		if err := handler(ctx, gwmux, s.gwConn); err != nil {
			s.Debugf("serve: %v\n", err)
			return err
		}
	}

	mux := http.NewServeMux()

	version := s.GetVersionPath()

	mux.Handle(filepath.Join("/", version, "/api"), gwmux)

	if s.swaggerFile != "" {
		mux.HandleFunc(filepath.Join("/", version, "/swagger.json"), serveSwaggerJSON(s.swaggerFile))
	}

	// Register the http server for gRPC gateway
	if certs, err := ParseCertificates(s.pubCert, s.privCert); err != nil {
		return err
	} else {
		s.httpServer = &http.Server{
			Addr:    endpoint,
			Handler: grpcHandlerFunc(s.grpcServer, mux),
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{certs},
				NextProtos:   []string{"h2"},
			},
		}
	}

	return nil
}

func serveSwaggerJSON(filepath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath)
	}
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
// Code from https://github.com/philips/grpc-gateway-example/blob/master/cmd/serve.go
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61

		// Checks if grpcServer is implemented
		if grpcServer != nil &&
			r.ProtoMajor == 2 &&
			strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}
