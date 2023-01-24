package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
// NOTE: this will fail if tls certs are not passed in and or it insecure is not passed
// GRPC Gateway and the GRPC Server cannot run on the same port if they are not tls
// the transport client cannot upgrade from http1 to http2 without it
func (s *Server) newGRPCGateway(ctx context.Context) error {

	dialOpts := []grpc.DialOption{}
	// Dial Credentials need to come from the outside if they are different than the local service
	// If no Certificates are pass we assume they are running on the same server

	if len(s.gatewayDialOptions) == 0 {
		// The cert should have been loaded in from the constructor, this is assuming that insecure wasn't passed in
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{Certificates: s.dialCerts})))
	}

	dialOpts = append(dialOpts, s.gatewayDialOptions...)

	// Register the gateway service handlers; service handlers currently only talk to same grpc.ClientConn
	// (joematpal) I do not see a use case where we need to have/maintain a one-to-many type of connection for the Gateway

	gwmux := runtime.NewServeMux(s.gatewayServerMuxOptions...)
	endpoint := s.getGatewayEndpoint()

	s.Debugf("gateway host: %v", endpoint)

	var err error
	s.gwConn, err = grpc.DialContext(ctx, endpoint, dialOpts...)
	if err != nil {
		return fmt.Errorf("dial context: %v", err)
	}

	s.Debugf("gw conn state: %v", s.gwConn.GetState())

	for _, handler := range s.gatewayServiceHandlers {
		if err := handler(ctx, gwmux, s.gwConn); err != nil {
			s.Debugf("register service handler: %v\n", err)
			return err
		}
	}

	var handler http.Handler

	// Swagger File Server
	if s.swaggerFile != "" {
		s.Debugf("serving swagger: %s", s.swaggerFile)
		// TODO: add the ability to change the swagger path?
		s.handler.Handle("/spec/v1/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, s.swaggerFile)
		}))
	}

	if s.handler == nil {
		handler = gwmux
	} else {
		s.handler.Handle("/", gwmux)
		handler = s.handler
	}
	fmt.Println("swagger folder", s.swaggerFolder)

	// Swagger Directory Server
	if s.swaggerFolder != "" {
		s.Debugf("serving swagger folder: %s", s.swaggerFile)
		// TODO: add the ability to change the swagger path?
		fmt.Println("swagger folder", s.swaggerFolder)
		s.handler.Handle("/spec/", http.StripPrefix("/spec/", http.FileServer(http.Dir(s.swaggerFolder))))
	}

	// if s.grpcServer != nil {
	// 	handler = s.grpcHandlerFunc(s.grpcServer, handler)
	// }

	// Register the http server for gRPC gateway
	if s.IsTLS() {
		s.Debugf("running tls")

		s.Debugf("insecureSkipVerify: %v", s.insecureSkipVerify)

		s.httpServer = &http.Server{
			Addr: fmt.Sprintf("%s:%s", s.host, s.port),
			TLSConfig: &tls.Config{
				Certificates:       s.dialCerts,
				ClientCAs:          s.getCertPool(),
				ClientAuth:         s.clientAuthType,
				InsecureSkipVerify: s.insecureSkipVerify,
				NextProtos:         []string{"h2"},
			},
			Handler: handler,
		}
	} else {
		s.httpServer = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", s.host, s.port),
			Handler: h2c.NewHandler(handler, &http2.Server{}),
		}
	}

	return nil
}

func (s *Server) getCertPool() *x509.CertPool {
	if s.clientCAs == nil {
		return x509.NewCertPool()
	}

	return s.clientCAs
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
// Code from https://github.com/philips/grpc-gateway-example/blob/master/cmd/serve.go
func (s *Server) grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		// Checks if grpcServer is implemented
		if grpcServer != nil &&
			r.ProtoMajor == 2 &&
			strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.Debugf("http2 path: %s", r.URL.Path)
			grpcServer.ServeHTTP(w, r)
		} else {
			s.Debugf("http path: %s", r.URL.Path)
			otherHandler.ServeHTTP(w, r)
		}
	})
}
