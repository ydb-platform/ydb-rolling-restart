package rolling

import (
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
)

var (
	ServiceOptionsMap = map[string]options.Options{}
	ServiceFactoryMap = map[string]ServiceFactory{}
)

type ServiceFactory func(o options.Options) (Service, error)
type Service interface {
	Prepare() error
	Filter(spec FilterNodeParams) []*Ydb_Maintenance.Node
	RestartNode(node *Ydb_Maintenance.Node) error
}

func init() {
	ServiceOptionsMap["mock"] = nil
	ServiceFactoryMap["mock"] = func(o options.Options) (Service, error) {
		return &mock{}, nil
	}
}

type mock struct{}

func (m *mock) Prepare() error { return nil }
func (m *mock) Filter(spec FilterNodeParams) []*Ydb_Maintenance.Node {
	// take only storage nodes
	return util.FilterBy(spec.AllNodes,
		func(node *Ydb_Maintenance.Node) bool {
			return node.GetStorage() != nil
		},
	)
}
func (m *mock) RestartNode(node *Ydb_Maintenance.Node) error { return nil }

type FilterNodeParams struct {
	Service         string
	AllTenants      []string
	AllNodes        []*Ydb_Maintenance.Node
	SelectedTenants []string
	SelectedNodeIds []string
}
