package flags

import "github.com/urfave/cli/v2"

func Join(flags ...[]cli.Flag) []cli.Flag {
	var out []cli.Flag
	for _, f := range flags {
		out = append(out, f...)
	}
	return out
}
