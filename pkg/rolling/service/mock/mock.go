package mock

import (
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/rolling/service"
)

const (
	name = "mock"
)

func init() {
	service.OptionsMap[name] = nil
	service.FactoryMap[name] = func(o options.Options) (service.Interface, error) {
		return &mock{}, nil
	}
}

type mock struct{}

func (m *mock) Prepare() error { return nil }
func (m *mock) Filter(spec service.FilterNodeParams) []*Ydb_Maintenance.Node {
	// take only storage nodes
	return util.FilterBy(spec.AllNodes,
		func(node *Ydb_Maintenance.Node) bool {
			return node.GetStorage() != nil
		},
	)
}
func (m *mock) RestartNode(*Ydb_Maintenance.Node) error { return nil }
