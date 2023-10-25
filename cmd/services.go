package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewServicesCommand(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "services",
		Short: "Fetch and output list of services of YDB Cluster",
		Long:  "Fetch and output list of services of YDB Cluster (long version)",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("services called")
		},
	}
}
