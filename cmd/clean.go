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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Validate(opts.GRPC, opts.CMS); err != nil {
				logger.Errorf("Failed to validate options: %+v", err)
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
			tasks, err := c.MaintenanceTasks()
			if err != nil {
				logger.Errorf("Failed to list tasks: %v", err)
				return nil
			}
			for _, taskId := range tasks {
				status, err := c.DropMaintenanceTask(taskId)
				if err != nil {
					logger.Errorf("Failed to drop maintenance task: %v", err)
					continue
				}
				logger.Infof("Drop task with id: %s status is %s", tasks, status)
			}
			return nil
		},
	}

	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())
	return cmd
}
