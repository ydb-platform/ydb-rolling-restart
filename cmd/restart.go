package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewRestartCommand(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Perform restart of YDB cluster",
		Long:  "Perform restart of YDB cluster (long version)",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("restart called")
		},
	}
}
