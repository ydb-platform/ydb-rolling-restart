package cms

import (
	"github.com/ydb-platform/ydb-go-genproto/Ydb_Cms_V1"
	"github.com/ydb-platform/ydb-go-genproto/protos/Ydb_Cms"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
)

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
	// todo: error handling, retries

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

	s := util.SortBy(o.Paths, func(l string, r string) bool { return l < r })
	return s, nil
}
