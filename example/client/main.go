package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	serverf "github.com/joematpal/go-server/flags"
	streamer "github.com/joematpal/go-server/pkg/streamer/v1"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:  "route_guide_client",
		Flags: serverf.GRPCFlags,
		Action: func(c *cli.Context) error {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			ctx := c.Context
			target := fmt.Sprintf("%s:%s", c.String(serverf.GRPCHost), c.String(serverf.GRPCPort))
			opts := []grpc.DialOption{
				grpc.WithInsecure(),
			}

			conn, err := grpc.DialContext(ctx, target, opts...)
			if err != nil {
				return fmt.Errorf("dial context: %v", err)
			}
			rg := streamer.NewStreamerServiceClient(conn)

			rrc, err := rg.StreamPoint(ctx)
			if err != nil {
				return fmt.Errorf("record route: %v", err)
			}

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()

				default:
					if err := rrc.Send(randomPoint(r)); err != nil {
						return fmt.Errorf("send: %v", err)
					}
				}
			}
		},
	}
}

func randomPoint(r *rand.Rand) *streamer.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &streamer.Point{Latitude: lat, Longitude: long}
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
