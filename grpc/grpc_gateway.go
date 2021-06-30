package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	assetfs "github.com/philips/go-bindata-assetfs"
	"github.com/philips/grpc-gateway-example/pkg/ui/data/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Register the http server for the GRPC Gateway.
// and Register the grpc client on the mux runtime (handler)
func (s *Server) newGRPCGateway() error {
	var creds credentials.TransportCredentials
	var err error
	certs, err := ParseCertificates(s.pubCert, s.privCert)
	if err != nil {
		return fmt.Errorf("parse certs: %v", err)
	}
	creds = Credentials(certs)

	// Register the gateway service handlers
	ctx := context.Background()
	gwmux := runtime.NewServeMux()
	endpoint := fmt.Sprintf("%s:%s", s.host, s.port)
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	for _, handler := range s.gatewayServiceHandlers {
		if err := handler(ctx, gwmux, endpoint, dopts); err != nil {
			s.Debugf("serve: %v\n", err)
			return err
		}
	}

	mux := http.NewServeMux()
	// TODO: Re-enable swagger
	// Create swagger endpoint
	// mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
	// 	io.Copy(w, strings.NewReader(pb.Swagger))
	// })

	mux.Handle("/", gwmux)
	// TODO: Re-enable swagger
	// serveSwagger(mux)

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

func serveSwagger(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
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
