package options

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

const (
	GRPCDefaultPort           = 2135
	GRPCDefaultTimeoutSeconds = 60
)

type GRPC struct {
	Addr           []string
	Port           int
	Secure         bool
	RootCA         string
	TimeoutSeconds int
}

func (grpc *GRPC) DefineFlags(fs *pflag.FlagSet) {
	fs.StringSliceVarP(&grpc.Addr, "grpc-address", "", grpc.Addr,
		"GRPC addresses which will be used to connect to cluster")
	fs.IntVarP(&grpc.Port, "grpc-port", "", GRPCDefaultPort,
		"GRPC port available on all addresses")
	fs.BoolVarP(&grpc.Secure, "grpc-secure", "", grpc.Secure,
		"GRPC or GRPCS protocol to use")
	fs.StringVarP(&grpc.RootCA, "grpc-secure-root-ca", "", grpc.RootCA,
		"GRPCS path to root CA")
	fs.IntVarP(&grpc.TimeoutSeconds, "grpc-api-timeout-seconds", "", GRPCDefaultTimeoutSeconds,
		"CMS API response timeout in seconds")
}

func (grpc *GRPC) Validate() error {
	if len(grpc.Addr) == 0 {
		return fmt.Errorf("specify not empty grpc addresses")
	}
	if grpc.Port < 0 || grpc.Port > 65536 {
		return fmt.Errorf("invalid port specified: %d, must be in range: (%d,%d)", grpc.Port, 1, 65536)
	}
	if grpc.RootCA != "" {
		if !grpc.Secure {
			return fmt.Errorf("root CA must be specified only for secure connection")
		}

		if _, err := os.Stat(grpc.RootCA); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("root CA file not found: %v", err)
		}
	}
	if grpc.TimeoutSeconds < 0 {
		return fmt.Errorf("invalid value specified: %d", grpc.TimeoutSeconds)
	}

	return nil
}
