/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRestartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Perform restart of YDB cluster",
		Long:  "Perform restart of YDB cluster (long version)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("restart called")
		},
	}
}
