/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewServicesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "services",
		Short: "Fetch and output list of services of YDB Cluster",
		Long:  "Fetch and output list of services of YDB Cluster (long version)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("services called")
		},
	}
}
