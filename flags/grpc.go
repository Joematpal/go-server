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
	GRPCHost            = "grpc-host"
	GRPCPort            = "grpc-port"
	GRPCTLS             = "grpc-tls"
	GRPCPubCert         = "grpc-pub-cert"
	GRPCPrivCert        = "grpc-priv-cert"
	GRPCGatewayHost     = "grpc-gw-host"
	GRPCGatewayPort     = "grpc-gw-port"
	GRPCGatewayPubCert  = "grpc-gw-pub-cert"
	GRPCGatewayPrivCert = "grpc-gw-priv-cert"
	SwaggerFile         = "swagger-file"
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
		Name:    GRPCGatewayHost,
		Usage:   "used to set the gateway host if it is different than the address for the local service",
		EnvVars: flagNamesToEnv(GRPCGatewayHost),
	},
	&cli.StringFlag{
		Name:    GRPCGatewayPort,
		Usage:   "used to set the gateway port if it is different than the address for the local service",
		EnvVars: flagNamesToEnv(GRPCGatewayPort),
	},
	&cli.StringFlag{
		Name:    GRPCGatewayPubCert,
		Usage:   "used to set the gateway certs if they are different than the ones for the local service",
		EnvVars: flagNamesToEnv(GRPCGatewayPubCert),
	},
	&cli.StringFlag{
		Name:    GRPCGatewayPrivCert,
		Usage:   "used to set the gateway certs if they are different than the ones for the local service",
		EnvVars: flagNamesToEnv(GRPCGatewayPrivCert),
	},
	&cli.StringFlag{
		Name:    SwaggerFile,
		Usage:   "set the filepath to `swagger.json`",
		EnvVars: flagNamesToEnv(SwaggerFile),
	},
}
