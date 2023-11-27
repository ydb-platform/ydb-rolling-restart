package cms

import (
	"context"

	"github.com/ydb-platform/ydb-go-genproto/Ydb_Cms_V1"
	"github.com/ydb-platform/ydb-go-genproto/draft/Ydb_Maintenance_V1"
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"
	"github.com/ydb-platform/ydb-go-genproto/protos/Ydb_Cms"
	"github.com/ydb-platform/ydb-go-genproto/protos/Ydb_Operations"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
)

type CMSClient struct {
	logger *zap.SugaredLogger
	f      *Factory
}

type operationResponse interface {
	GetOperation() *Ydb_Operations.Operation
}

func NewCMSClient(logger *zap.SugaredLogger, f *Factory) *CMSClient {
	return &CMSClient{
		logger: logger,
		f:      f,
	}
}

func (c *CMSClient) Tenants() ([]string, error) {
	result := Ydb_Cms.ListDatabasesResult{}
	_, err := c.ExecuteCMSMethod(&result, func(ctx context.Context, cl Ydb_Cms_V1.CmsServiceClient) (operationResponse, error) {
		return cl.ListDatabases(ctx, &Ydb_Cms.ListDatabasesRequest{OperationParams: c.f.OperationParams()})
	})
	if err != nil {
		return nil, err
	}

	s := util.SortBy(result.Paths,
		func(l string, r string) bool {
			return l < r
		},
	)
	return s, nil
}

func (c *CMSClient) Nodes() ([]*Ydb_Maintenance.Node, error) {
	result := Ydb_Maintenance.ListClusterNodesResult{}
	_, err := c.ExecuteMaintenanceMethod(&result,
		func(ctx context.Context, cl Ydb_Maintenance_V1.MaintenanceServiceClient) (operationResponse, error) {
			return cl.ListClusterNodes(ctx, &Ydb_Maintenance.ListClusterNodesRequest{OperationParams: c.f.OperationParams()})
		},
	)
	if err != nil {
		return nil, err
	}

	nodes := util.SortBy(result.Nodes,
		func(l *Ydb_Maintenance.Node, r *Ydb_Maintenance.Node) bool {
			return l.NodeId < r.NodeId
		},
	)

	return nodes, nil
}

func (c *CMSClient) MaintenanceTasks() ([]string, error) {
	result := Ydb_Maintenance.ListMaintenanceTasksResult{}
	_, err := c.ExecuteMaintenanceMethod(&result,
		func(ctx context.Context, cl Ydb_Maintenance_V1.MaintenanceServiceClient) (operationResponse, error) {
			return cl.ListMaintenanceTasks(ctx,
				&Ydb_Maintenance.ListMaintenanceTasksRequest{
					OperationParams: c.f.OperationParams(),
					User:            util.Pointer(c.f.User()),
				},
			)
		},
	)
	if err != nil {
		return nil, err
	}

	return result.TasksUids, nil
}

func (c *CMSClient) DropMaintenanceTask(taskId string) (string, error) {
	op, err := c.ExecuteMaintenanceMethod(nil,
		func(ctx context.Context, cl Ydb_Maintenance_V1.MaintenanceServiceClient) (operationResponse, error) {
			return cl.DropMaintenanceTask(ctx, &Ydb_Maintenance.DropMaintenanceTaskRequest{
				OperationParams: c.f.OperationParams(),
				TaskUid:         taskId,
			})
		},
	)
	if err != nil {
		return "", err
	}

	return op.Status.String(), nil
}

func (c *CMSClient) ExecuteMaintenanceMethod(
	out proto.Message,
	method func(context.Context, Ydb_Maintenance_V1.MaintenanceServiceClient) (operationResponse, error),
) (*Ydb_Operations.Operation, error) {
	// todo:
	// 	 1, error handling ??
	//   2. retries
	//   3. operation result check

	cc, err := c.f.Connection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := c.f.Context()
	defer cancel()

	cl := Ydb_Maintenance_V1.NewMaintenanceServiceClient(cc)
	r, err := method(ctx, cl)
	if err != nil {
		return nil, err
	}
	op := r.GetOperation()

	if out == nil {
		return op, nil
	}

	if err := op.Result.UnmarshalTo(out); err != nil {
		return op, err
	}

	return op, nil
}

func (c *CMSClient) ExecuteCMSMethod(
	out proto.Message,
	method func(context.Context, Ydb_Cms_V1.CmsServiceClient) (operationResponse, error),
) (*Ydb_Operations.Operation, error) {
	cc, err := c.f.Connection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := c.f.Context()
	defer cancel()

	cl := Ydb_Cms_V1.NewCmsServiceClient(cc)
	r, err := method(ctx, cl)
	if err != nil {
		return nil, err
	}
	op := r.GetOperation()

	if out == nil {
		return op, nil
	}

	if err := op.Result.UnmarshalTo(out); err != nil {
		return op, err
	}

	return op, nil
}
