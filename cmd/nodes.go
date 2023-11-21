package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/cms"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
)

type NodesOptions struct {
	CMS  *options.CMS
	GRPC *options.GRPC
}

func NewNodesCommand(lf *zap.Logger) *cobra.Command {
	logger := lf.Sugar()
	opts := TenantsOptions{
		CMS:  &options.CMS{},
		GRPC: &options.GRPC{},
	}
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Fetch and output list of nodes of YDB Cluster",
		Long:  "Fetch and output list of nodes of YDB Cluster (long version)",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			nodes, err := c.Nodes()
			if err != nil {
				logger.Errorf("Failed to list nodes: %v", err)
				return err
			}

			msg := util.Join(nodes, "\n",
				func(node *Ydb_Maintenance.ListClusterNodesResponse_Node) string {
					return fmt.Sprintf("%d\t%s\t%s", node.NodeId, node.Fqdn, node.State)
				},
			)
			logger.Infof("Nodes:\n%s", msg)

			return nil
		},
	}

	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())
	return cmd
}
