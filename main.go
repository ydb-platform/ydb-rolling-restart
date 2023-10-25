/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ydb-platform/ydb-rolling-restart/cmd"
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func createLogger(level string) (zap.AtomicLevel, *zap.Logger) {
	atom, _ := zap.ParseAtomicLevel(level)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return atom, zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
}

func main() {
	logLevel := "info"
	logLevelSetter, logger := createLogger(logLevel)
	root := &cobra.Command{
		Use:   "ydb-rolling-restart",
		Short: "The rolling restart utility",
		Long:  "The rolling restart utility (long version)",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			lvc, err := zapcore.ParseLevel(logLevel)
			if err != nil {
				logger.Warn("Failed to set level")
				return
			}
			logLevelSetter.SetLevel(lvc)
		},
	}
	defer logger.Sync()

	root.PersistentFlags().StringVarP(&logLevel, "log-level", "", logLevel, "Logging level")
	root.AddCommand(
		cmd.NewServicesCommand(logger),
		cmd.NewTenantsCommand(logger),
		cmd.NewRestartCommand(logger),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
