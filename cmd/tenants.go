package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ydb-platform/ydb-go-genproto/Ydb_Cms_V1"
	"github.com/ydb-platform/ydb-go-genproto/protos/Ydb_Cms"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/cms"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
)

type TenantsOptions struct {
	CMS  *options.CMS
	GRPC *options.GRPC
}

func NewTenantsCommand(lf *zap.Logger) *cobra.Command {
	logger := lf.Sugar()
	opts := TenantsOptions{
		CMS:  &options.CMS{},
		GRPC: &options.GRPC{},
	}
	cmd := &cobra.Command{
		Use:   "tenants",
		Short: "Fetch and output list of tenants of YDB Cluster",
		Long:  "Fetch and output list of tenants of YDB Cluster (long version)",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := options.Validate(opts.GRPC, opts.CMS); err != nil {
				logger.Error("Failed to validate options", zap.Error(err))
				return err
			}

			logger.Info("Started")
			f := cms.NewConnectionFactory(
				*opts.GRPC,
				opts.CMS.Auth,
				opts.CMS.User,
			)
			logger.Info("Create connection")
			cc, err := f.Connection()
			if err != nil {
				logger.Errorf("%+v", err)
				return err
			}

			cl := Ydb_Cms_V1.NewCmsServiceClient(cc)
			ctx, cancel := f.Context()
			defer cancel()

			logger.Info("Invoke list databases")
			r, err := cl.ListDatabases(ctx,
				&Ydb_Cms.ListDatabasesRequest{
					OperationParams: f.OperationParams(),
				},
			)

			if err != nil {
				logger.Errorf("%+v", err)
				return err
			}
			logger.Infof("%+v", r)

			return nil
		},
	}
	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())

	return cmd
}
