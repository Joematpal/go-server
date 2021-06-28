package main

import (
	"fmt"

	"github.com/digital-dream-labs/go-server/flags"
	"github.com/digital-dream-labs/go-server/grpc"
	cli "github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:  "example grpc server",
		Flags: flags.GRPCFlags,
		Action: func(c *cli.Context) error {
			opts := []grpc.Option{}
			srv, err := grpc.New(opts...)
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

}
