package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/cms"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
)

type CleanOptions struct {
	CMS  *options.CMS
	GRPC *options.GRPC
}

func NewCleanCommand(lf *zap.Logger) *cobra.Command {
	logger := lf.Sugar()
	opts := CleanOptions{
		CMS:  &options.CMS{},
		GRPC: &options.GRPC{},
	}
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Perform cleanup of management requests in cluster",
		Long:  "Perform cleanup of management requests in cluster (long version)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return options.Validate(opts.GRPC, opts.CMS)
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := cms.NewCMSClient(
				logger,
				cms.NewConnectionFactory(
					*opts.GRPC,
					opts.CMS.Auth,
					opts.CMS.User,
				),
			)
			tasks, err := c.MaintenanceTasks()
			if err != nil {
				logger.Errorf("Failed to list tasks: %v", err)
				return
			}
			for _, taskId := range tasks {
				status, err := c.DropMaintenanceTask(taskId)
				if err != nil {
					logger.Errorf("Failed to drop maintenance task: %v", err)
					continue
				}
				logger.Infof("Drop task with id: %s status is %s", tasks, status)
			}
		},
	}

	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())
	return cmd
}
