package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/cms"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/rolling"
	_ "github.com/ydb-platform/ydb-rolling-restart/pkg/rolling/service/mock"
	_ "github.com/ydb-platform/ydb-rolling-restart/pkg/rolling/service/ssh"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return options.Validate(opts.GRPC, opts.CMS, opts.Rolling)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// todo: any cleanup required?
			var err error

			client := cms.NewCMSClient(logger,
				cms.NewConnectionFactory(
					*opts.GRPC,
					opts.CMS.Auth,
					opts.CMS.User,
				),
			)
			r := rolling.New(client, logger, opts.Rolling)

			if opts.Continue {
				logger.Info("Continue previous rolling restart")
				err = r.RestartPrevious()
			} else {
				logger.Info("Start rolling restart")
				err = r.Restart()
			}

			if err != nil {
				logger.Errorf("Failed to complete restart: %+v", err)
			} else {
				logger.Info("Restart completed successfully")
			}
		},
	}

	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())
	opts.Rolling.DefineFlags(cmd.Flags())
	cmd.Flags().BoolVarP(&opts.Continue, "continue", "", opts.Continue,
		"Continue previous restart")
	return cmd
}
