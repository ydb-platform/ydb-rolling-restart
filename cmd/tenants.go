package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
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

			c := cms.NewCMSClient(
				logger,
				cms.NewConnectionFactory(
					*opts.GRPC,
					opts.CMS.Auth,
					opts.CMS.User,
				),
			)
			tenants, err := c.Tenants()
			if err != nil {
				logger.Errorf("Failed to list tenants: %v", err)
				return err
			}

			msg := util.Join(tenants, "\n", func(s string) string { return s })
			logger.Infof("Tenants:\n%s", msg)

			return nil
		},
	}
	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())

	return cmd
}
