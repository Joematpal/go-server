package main

import (
	"fmt"

	"github.com/digital-dream-labs/go-server/flags"
	grpcp "github.com/digital-dream-labs/go-server/grpc"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:  "example grpc server",
		Flags: flags.GRPCFlags,
		Action: func(c *cli.Context) error {
			opts := []grpcp.Option{
				grpcp.WithRegisterService(func(*grpc.Server) {
					// fmt.Println("Add your service here")
				}),
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

}
