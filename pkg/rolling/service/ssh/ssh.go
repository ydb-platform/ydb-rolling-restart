package ssh

import (
	"fmt"

	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/rolling/service"
)

const (
	name = "ssh"
)

func init() {
	service.FactoryMap[name] = New
}

type ssh struct {
	options *opts
}

func New(o options.Options) (service.Interface, error) {
	opts, ok := o.(*opts)
	if !ok {
		return nil, fmt.Errorf("invalid options type specified: %+v", o)
	}
	return &ssh{
		options: opts,
	}, nil
}

func (s *ssh) Prepare() error { return nil }

func (s *ssh) Filter(spec service.FilterNodeParams) []*Ydb_Maintenance.Node {
	nodes := util.FilterBy(spec.AllNodes,
		func(node *Ydb_Maintenance.Node) bool {
			if s.options.ServiceName == ServiceTypeStorage {
				return node.GetStorage() != nil
			}
			return node.GetDynamic() != nil
		},
	)

	if s.options.ServiceType == ServiceTypeDynamic && len(spec.SelectedTenants) > 0 {
		nodes = util.FilterBy(nodes,
			func(node *Ydb_Maintenance.Node) bool {
				return util.Contains(spec.SelectedTenants, node.GetDynamic().Tenant)
			},
		)
	}

	if len(spec.SelectedNodeIds) > 0 {
		nodes = util.FilterBy(nodes,
			func(node *Ydb_Maintenance.Node) bool {
				return util.Contains(spec.SelectedNodeIds, node.NodeId)
			},
		)
	}

	return nodes
}

func (s *ssh) RestartNode(node *Ydb_Maintenance.Node) error {
	return nil
}
