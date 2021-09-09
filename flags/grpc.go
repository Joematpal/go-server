package flags

import (
	"strings"

	"github.com/urfave/cli/v2"
)

// flagNamesToEnv converts flags to a ENV format
func flagNamesToEnv(names ...string) []string {
	out := []string{}
	for _, name := range names {
		out = append(out, flagNameToEnv(name))
	}
	return out
}

// flagNameToEnv converts a flag to an ENV format
func flagNameToEnv(name string) string {
	return strings.ReplaceAll(strings.ToUpper(name), "-", "_")
}

var (
	GRPCHost                  = "grpc-host"
	GRPCPort                  = "grpc-port"
	GRPCTLS                   = "grpc-tls"
	GRPCPubCert               = "grpc-pub-cert"
	GRPCPrivCert              = "grpc-priv-cert"
	GRPCGatewayClientHost     = "grpc-gwc-host"
	GRPCGatewayClientPort     = "grpc-gwc-port"
	GRPCGatewayClientPubCert  = "grpc-gwc-pub-cert"
	GRPCGatewayClientPrivCert = "grpc-gwc-priv-cert"
	SwaggerFile               = "swagger-file"
)

var GRPCFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    GRPCHost,
		EnvVars: flagNamesToEnv(GRPCHost),
	},
	&cli.StringFlag{
		Name:    GRPCPort,
		EnvVars: flagNamesToEnv(GRPCPort),
	},
	&cli.BoolFlag{
		Name:    GRPCTLS,
		EnvVars: flagNamesToEnv(GRPCTLS),
	},
	&cli.StringFlag{
		Name:    GRPCPubCert,
		EnvVars: flagNamesToEnv(GRPCPubCert),
	},
	&cli.StringFlag{
		Name:    GRPCPrivCert,
		EnvVars: flagNamesToEnv(GRPCPrivCert),
	},
	&cli.StringFlag{
		Name:    GRPCGatewayClientHost,
		Usage:   "used to set the gateway client host if different than local service",
		EnvVars: flagNamesToEnv(GRPCGatewayClientHost),
	},
	&cli.StringFlag{
		Name:    GRPCGatewayClientPort,
		Usage:   "used to set the gateway client port if different than local service",
		EnvVars: flagNamesToEnv(GRPCGatewayClientPort),
	},
	&cli.StringFlag{
		Name:    GRPCGatewayClientPubCert,
		Usage:   "used to set the gateway certs if they are different than the ones for the local service",
		EnvVars: flagNamesToEnv(GRPCGatewayClientPubCert),
	},
	&cli.StringFlag{
		Name:    GRPCGatewayClientPrivCert,
		Usage:   "used to set the gateway certs if they are different than the ones for the local service",
		EnvVars: flagNamesToEnv(GRPCGatewayClientPrivCert),
	},
	&cli.StringFlag{
		Name:    SwaggerFile,
		Usage:   "set the filepath to `swagger.json`",
		EnvVars: flagNamesToEnv(SwaggerFile),
	},
}
