package options

import (
	"github.com/ydb-platform/ydb-rolling-restart/internal/util"

	"github.com/spf13/pflag"
	"go.uber.org/multierr"
)

// Options is an interface to defile options flags and validation logic
type Options interface {
	DefineFlags(fs *pflag.FlagSet)
	Validate() error
}

func Validate(options ...Options) error {
	return multierr.Combine(
		util.Convert(options, func(o Options) error { return o.Validate() })...,
	)
}
