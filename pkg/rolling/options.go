package rolling

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
)

type Options struct {
	ExecutorType string
}

func (o *Options) DefineFlags(fs *pflag.FlagSet) {
	executors := util.Join(util.Keys(ExecutorOptionsMap), ", ",
		func(s string) string {
			return s
		},
	)

	fs.StringVarP(&o.ExecutorType, "executor-type", "", o.ExecutorType,
		fmt.Sprintf("Executor type. Available choices: %s", executors))

	for _, executor := range ExecutorOptionsMap {
		executor.DefineFlags(fs)
	}
}

func (o *Options) Validate() error {
	eo, exists := ExecutorOptionsMap[o.ExecutorType]
	if !exists {
		return fmt.Errorf("specified not supported executor: %s", o.ExecutorType)
	}
	return eo.Validate()
}
