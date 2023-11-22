package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
)

type RestartOptions struct {
	CMS  *options.CMS
	GRPC *options.GRPC
}

func NewRestartCommand(lf *zap.Logger) *cobra.Command {
	logger := lf.Sugar()
	opts := RestartOptions{
		CMS:  &options.CMS{},
		GRPC: &options.GRPC{},
	}
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Perform restart of YDB cluster",
		Long:  "Perform restart of YDB cluster (long version)",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("restart called")
			return nil
		},
	}

	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())
	return cmd
}
