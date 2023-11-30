package rolling

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
)

var (
	AvailabilityModes = []string{"strong", "weak", "force"}
)

type Options struct {
	Mode    string
	Service string
	Tenants []string
	Nodes   []string
}

func (o *Options) DefineFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Service, "service", "", o.Service,
		fmt.Sprintf("Service type. Available choices: %s", util.JoinStrings(util.Keys(ServiceOptionsMap), ", ")))

	fs.StringVarP(&o.Mode, "availability-mode", "", AvailabilityModes[0],
		fmt.Sprintf("Availability mode. Available choices: %s", util.JoinStrings(AvailabilityModes, ", ")))

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

	if !util.Contains(AvailabilityModes, o.Mode) {
		return fmt.Errorf("specified not supported availability mode: %s", o.AvailabilityMode)
	}
	return so.Validate()
}

func (o *Options) AvailabilityMode() Ydb_Maintenance.AvailabilityMode {
	title := strings.ToUpper(fmt.Sprintf("availability_mode_%s", o.Mode))
	value := Ydb_Maintenance.AvailabilityMode_value[title]

	return Ydb_Maintenance.AvailabilityMode(value)
}
