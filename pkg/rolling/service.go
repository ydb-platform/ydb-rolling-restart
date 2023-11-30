package rolling

import (
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"
)

var (
	ServiceOptionsMap = map[string]Options{}
	ServiceFactoryMap = map[string]ServiceFactory{}
)

type ServiceFactory func(o Options) (Service, error)
type Service interface {
	Prepare() error
	Filter(spec FilterNodeParams) []*Ydb_Maintenance.Node
	RestartNode() error
}

type FilterNodeParams struct {
	Service         string
	AllTenants      []string
	AllNodes        []*Ydb_Maintenance.Node
	SelectedTenants []string
	SelectedNodeIds []string
}
