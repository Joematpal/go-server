package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/digital-dream-labs/go-logger"
	loggerf "github.com/digital-dream-labs/go-logger/flags"

	"github.com/digital-dream-labs/go-server/example/server/internal/flags"
	route_guide "github.com/digital-dream-labs/go-server/example/server/internal/route_guide"
	serverf "github.com/digital-dream-labs/go-server/flags"
	grpcp "github.com/digital-dream-labs/go-server/grpc"
	streamer "github.com/digital-dream-labs/go-server/pkg/streamer/v1"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:  "example grpc server",
		Flags: flags.Join(serverf.GRPCFlags, loggerf.LogFlags),
		Action: func(c *cli.Context) error {
			logOpts := []logger.Option{
				logger.WithEnv(logger.LogEnv(c.Int(loggerf.LogEnv))),
				logger.WithLevel(logger.LogLevel(c.Int(loggerf.LogLevel))),
				logger.WithLogStacktrace(c.Bool(loggerf.LogStacktrace)),
			}
			logr, err := logger.New(logOpts...)
			if err != nil {
				return fmt.Errorf("new logger: %v", err)
			}

			opts := []grpcp.Option{
				grpcp.WithRegisterService(func(s *grpc.Server) {
					rg := route_guide.New(logr)
					streamer.RegisterStreamerServer(s, rg)
				}),
				grpcp.WithGatewayServiceHandlers(streamer.RegisterStreamerHandler),
				grpcp.WithLogger(logr),
			}

			if host := c.String(serverf.GRPCHost); host != "" {
				opts = append(opts, grpcp.WithHost(host))
			}

			if port := c.String(serverf.GRPCPort); port != "" {
				opts = append(opts, grpcp.WithPort(port))
			}

			if tls := c.Bool(serverf.GRPCTLS); tls {
				opts = append(opts, grpcp.WithTLS(tls))
			}
			if pubCert := c.String(serverf.GRPCPubCert); pubCert != "" {
				opts = append(opts, grpcp.WithPubCert(pubCert))
			}
			if privCert := c.String(serverf.GRPCPrivCert); privCert != "" {
				opts = append(opts, grpcp.WithPrivCert(privCert))
			}

			srv, err := grpcp.New(opts...)
			if err != nil {
				return err
			}

			if err := srv.StartWithContext(c.Context); err != nil {
				return fmt.Errorf("start: %v", err)
			}
			return nil
		},
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
	}()
	if err := NewApp().RunContext(ctx, os.Args); err != nil {
		panic(err)
	}
}
