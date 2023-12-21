package service

import (
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
)

var (
	OptionsMap = map[string]options.Options{}
	FactoryMap = map[string]Factory{}
)

type (
	Factory   func(o options.Options) (Interface, error)
	Interface interface {
		Prepare() error
		Filter(spec FilterNodeParams) []*Ydb_Maintenance.Node
		RestartNode(node *Ydb_Maintenance.Node) error
	}
)

type FilterNodeParams struct {
	Service         string
	AllTenants      []string
	AllNodes        []*Ydb_Maintenance.Node
	SelectedTenants []string
	SelectedNodeIds []uint32
}
