package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type TenantsOptions struct {
}

func NewTenantsCommand(logger *zap.Logger) *cobra.Command {
	c := &cobra.Command{
		Use:   "tenants",
		Short: "Fetch and output list of tenants of YDB Cluster",
		Long:  "Fetch and output list of tenants of YDB Cluster (long version)",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("tenants called")
		},
	}

	return c
}
