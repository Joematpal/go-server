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
	GRPCHost     = "grpc-host"
	GRPCPort     = "grpc-port"
	GRPCTLS      = "grpc-tls"
	GRPCPubCert  = "grpc-pub-cert"
	GRPCPrivCert = "grpc-priv-cert"
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
}
