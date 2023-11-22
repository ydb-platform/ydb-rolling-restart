package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
)

type CleanOptions struct {
	CMS  *options.CMS
	GRPC *options.GRPC
}

func NewCleanCommand(lf *zap.Logger) *cobra.Command {
	logger := lf.Sugar()
	opts := RestartOptions{
		CMS:  &options.CMS{},
		GRPC: &options.GRPC{},
	}
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Perform cleanup of management requests in cluster",
		Long:  "Perform cleanup of management requests in cluster (long version)",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("clean called")
			return nil
		},
	}

	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())
	return cmd
}
