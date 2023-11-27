package rolling

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
)

type Options struct {
	Service string
	Tenants []string
	Nodes   []string
}

func (o *Options) DefineFlags(fs *pflag.FlagSet) {
	services := util.Join(util.Keys(ServiceOptionsMap), ", ",
		func(s string) string {
			return s
		},
	)

	fs.StringVarP(&o.Service, "service", "", o.Service,
		fmt.Sprintf("Service type. Available choices: %s", services))

	fs.StringArrayVarP(&o.Tenants, "tenants", "", o.Tenants,
		"Restart only specified tenants")

	fs.StringArrayVarP(&o.Nodes, "nodes", "", o.Nodes,
		"Restart only specified nodes")

	for _, executor := range ServiceOptionsMap {
		executor.DefineFlags(fs)
	}
}

func (o *Options) Validate() error {
	so, exists := ServiceOptionsMap[o.Service]
	if !exists {
		return fmt.Errorf("specified not supported service: %s", o.Service)
	}
	return so.Validate()
}
