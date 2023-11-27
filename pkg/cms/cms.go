package cms

import (
	"github.com/ydb-platform/ydb-go-genproto/Ydb_Cms_V1"
	"github.com/ydb-platform/ydb-go-genproto/draft/Ydb_Maintenance_V1"
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"
	"github.com/ydb-platform/ydb-go-genproto/protos/Ydb_Cms"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
)

// todo: add error handling, retries
// todo: add operation result check

type CMSClient struct {
	logger *zap.SugaredLogger
	f      *Factory
}

func NewCMSClient(logger *zap.SugaredLogger, f *Factory) *CMSClient {
	return &CMSClient{
		logger: logger,
		f:      f,
	}
}

func (c *CMSClient) Tenants() ([]string, error) {
	cc, err := c.f.Connection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := c.f.Context()
	defer cancel()

	cl := Ydb_Cms_V1.NewCmsServiceClient(cc)
	r, err := cl.ListDatabases(ctx,
		&Ydb_Cms.ListDatabasesRequest{
			OperationParams: c.f.OperationParams(),
		},
	)
	if err != nil {
		return nil, err
	}

	o := Ydb_Cms.ListDatabasesResult{}
	if err := r.Operation.Result.UnmarshalTo(&o); err != nil {
		return nil, err
	}

	s := util.SortBy(o.Paths,
		func(l string, r string) bool {
			return l < r
		},
	)
	return s, nil
}

func (c *CMSClient) Nodes() ([]*Ydb_Maintenance.Node, error) {
	cc, err := c.f.Connection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := c.f.Context()
	defer cancel()

	cl := Ydb_Maintenance_V1.NewMaintenanceServiceClient(cc)
	r, err := cl.ListClusterNodes(ctx,
		&Ydb_Maintenance.ListClusterNodesRequest{},
	)
	if err != nil {
		return nil, err
	}

	o := Ydb_Maintenance.ListClusterNodesResult{}
	if err := r.Operation.Result.UnmarshalTo(&o); err != nil {
		return nil, err
	}

	s := util.SortBy(o.Nodes,
		func(l *Ydb_Maintenance.Node, r *Ydb_Maintenance.Node) bool {
			return l.NodeId < r.NodeId
		},
	)
	return s, nil
}

func (c *CMSClient) MaintenanceTasks() ([]string, error) {
	cc, err := c.f.Connection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := c.f.Context()
	defer cancel()

	cl := Ydb_Maintenance_V1.NewMaintenanceServiceClient(cc)
	r, err := cl.ListMaintenanceTasks(ctx, &Ydb_Maintenance.ListMaintenanceTasksRequest{
		User: util.Pointer(c.f.User()),
	})
	if err != nil {
		return nil, err
	}

	o := Ydb_Maintenance.ListMaintenanceTasksResult{}
	if err := r.Operation.Result.UnmarshalTo(&o); err != nil {
		return nil, err
	}

	return o.TasksUids, nil
}

func (c *CMSClient) DropMaintenanceTask(taskId string) (string, error) {
	cc, err := c.f.Connection()
	if err != nil {
		return "", err
	}

	ctx, cancel := c.f.Context()
	defer cancel()

	cl := Ydb_Maintenance_V1.NewMaintenanceServiceClient(cc)
	r, err := cl.DropMaintenanceTask(ctx, &Ydb_Maintenance.DropMaintenanceTaskRequest{TaskUid: taskId})
	if err != nil {
		return "", err
	}

	return r.Operation.Id, nil
}
