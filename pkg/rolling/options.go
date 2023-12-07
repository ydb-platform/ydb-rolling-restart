package rolling

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/rolling/service"
)

var (
	AvailabilityModes = []string{"strong", "weak", "force"}
)

type Options struct {
	AvailabilityMode   string
	Service            string
	Tenants            []string
	Nodes              []string
	RestartDuration    int
	RestartRetryNumber int
}

func (o *Options) DefineFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Service, "service", "", o.Service,
		fmt.Sprintf("Service type. Available choices: %s", util.JoinStrings(util.Keys(service.OptionsMap), ", ")))

	fs.StringVarP(&o.AvailabilityMode, "availability-mode", "", AvailabilityModes[0],
		fmt.Sprintf("Availability mode. Available choices: %s", util.JoinStrings(AvailabilityModes, ", ")))

	fs.IntVarP(&o.RestartDuration, "restart-duration", "", 60,
		"Restart duration in seconds")

	fs.IntVarP(&o.RestartRetryNumber, "restart-retry-number", "", 3,
		"Retry number of restart")

	fs.StringArrayVarP(&o.Tenants, "tenants", "", o.Tenants,
		"Restart only specified tenants")

	fs.StringArrayVarP(&o.Nodes, "nodes", "", o.Nodes,
		"Restart only specified nodes")

	for _, opts := range service.OptionsMap {
		if opts != nil {
			opts.DefineFlags(fs)
		}
	}
}

func (o *Options) Validate() error {
	opts, exists := service.OptionsMap[o.Service]
	if !exists {
		return fmt.Errorf("specified not supported service: %s", o.Service)
	}

	if !util.Contains(AvailabilityModes, o.AvailabilityMode) {
		return fmt.Errorf("specified not supported availability mode: %s", o.AvailabilityMode)
	}

	if o.RestartDuration < 0 {
		return fmt.Errorf("specified invalid restart duration seconds: %d. Must be positive", o.RestartDuration)
	}

	if o.RestartRetryNumber < 0 {
		return fmt.Errorf("specified invalid restart retry number: %d. Must be positive", o.RestartRetryNumber)
	}

	if _, err := o.GetNodeIds(); err != nil {
		return err
	}

	if opts != nil {
		return opts.Validate()
	}
	return nil
}

func (o *Options) GetAvailabilityMode() Ydb_Maintenance.AvailabilityMode {
	title := strings.ToUpper(fmt.Sprintf("availability_mode_%s", o.AvailabilityMode))
	value := Ydb_Maintenance.AvailabilityMode_value[title]

	return Ydb_Maintenance.AvailabilityMode(value)
}

func (o *Options) GetRestartDuration() *durationpb.Duration {
	return durationpb.New(time.Second * time.Duration(o.RestartDuration) * time.Duration(o.RestartRetryNumber))
}

func (o *Options) GetNodeIds() ([]uint32, error) {
	ids := make([]uint32, 0, len(o.Nodes))

	for _, nodeId := range o.Nodes {
		id, err := strconv.Atoi(nodeId)
		if err != nil {
			return nil, fmt.Errorf("failed to parse node id: %+v", err)
		}
		if id < 0 {
			return nil, fmt.Errorf("invalid node id specified: %d, must be positive", id)
		}
		ids = append(ids, uint32(id))
	}

	return ids, nil
}
