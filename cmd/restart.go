package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/cms"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/rolling"
)

type RestartOptions struct {
	CMS      *options.CMS
	GRPC     *options.GRPC
	Rolling  *rolling.Options
	Continue bool
}

func NewRestartCommand(lf *zap.Logger) *cobra.Command {
	logger := lf.Sugar()
	opts := RestartOptions{
		CMS:     &options.CMS{},
		GRPC:    &options.GRPC{},
		Rolling: &rolling.Options{},
	}
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Perform restart of YDB cluster",
		Long:  "Perform restart of YDB cluster (long version)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Validate(opts.GRPC, opts.CMS, opts.Rolling); err != nil {
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
			r := rolling.New(c, logger)

			// todo: any cleanup?
			// todo: logging here

			if opts.Continue {
				return r.RestartPrevious()
			}

			return r.Restart()
		},
	}

	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())
	opts.Rolling.DefineFlags(cmd.Flags())
	return cmd
}
