package rolling

import (
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/cms"
)

type Rolling struct {
	logger *zap.SugaredLogger
	cms    *cms.CMSClient
	opts   *Options
}

func New(cms *cms.CMSClient, logger *zap.SugaredLogger, opts *Options) *Rolling {
	return &Rolling{
		cms:    cms,
		logger: logger,
		opts:   opts,
	}
}

func (r *Rolling) Restart() error {
	factory := ExecutorFactoryMap[r.opts.ExecutorType]
	ex, err := factory(ExecutorOptionsMap[r.opts.ExecutorType])
	if err != nil {
		return err
	}

	if err := ex.Prepare(); err != nil {
		return err
	}

	// todo:
	//  1. filter out nodes by tenant/specified nodes
	//  2. invoke CreateMaintenanceTask
	//  3. run loop

	return nil
}

func (r *Rolling) RestartPrevious() error {

	// todo:
	//  1. find previous restart
	//  2. run loop

	return nil
}

func (r *Rolling) loop() {
	for r.next() {
		// todo: sleep if required
	}
}

func (r *Rolling) next() bool {
	return true
}