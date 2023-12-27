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
			nodes, err := c.Nodes()
			if err != nil {
				logger.Errorf("Failed to list nodes: %v", err)
				return
			}

			msg := util.Join(nodes, "\n",
				func(node *Ydb_Maintenance.Node) string {
					return fmt.Sprintf("%d\t%s\t%s", node.NodeId, node.Host, node.State)
				},
			)
			logger.Infof("Nodes:\n%s", msg)
		},
	}

	opts.CMS.DefineFlags(cmd.Flags())
	opts.GRPC.DefineFlags(cmd.Flags())
	return cmd
}
