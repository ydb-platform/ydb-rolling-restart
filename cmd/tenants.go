/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewTenantsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tenants",
		Short: "Fetch and output list of tenants of YDB Cluster",
		Long:  "Fetch and output list of tenants of YDB Cluster (long version)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("tenants called")
		},
	}
}
