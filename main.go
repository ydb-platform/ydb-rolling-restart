/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ydb-platform/ydb-rolling-restart/cmd"
)

func main() {
	var root = &cobra.Command{
		Use:   "ydb-rolling-restart",
		Short: "The rolling restart utility",
		Long:  "The rolling restart utility (long version)",
	}

	root.AddCommand(
		cmd.NewServicesCommand(),
		cmd.NewTenantsCommand(),
		cmd.NewRestartCommand(),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
